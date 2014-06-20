package csvdatafile

import (
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego/orm"
	"os"
	"strings"
)

type CSVDataFile struct {
	Table              string
	File               string
	fh                 *os.File
	writer             *csv.Writer
	hasImportError     bool
	FieldsTerminatedBy rune
	EnclosedBy         string
	LinesTerminatedBy  string
	Fields             []string
}

func New(table, file string) *CSVDataFile {
	csvDataFile := new(CSVDataFile)
	csvDataFile.Table = table
	csvDataFile.File = file
	csvDataFile.FieldsTerminatedBy = ';'
	return csvDataFile
}

func (c *CSVDataFile) AddField(field string) {
	c.Fields = append(c.Fields, field)
}

func (c *CSVDataFile) Create() error {
	var err error
	c.fh, err = os.Create(c.File)
	return err
}

func (c *CSVDataFile) Open() error {
	var err error
	c.fh, err = os.Open(c.File)
	return err
}

func (c *CSVDataFile) AppendRow(record ...interface{}) error {
	var err error
	c.writer = csv.NewWriter(c.fh)
	c.writer.Comma = c.FieldsTerminatedBy
	recordToWrite := make([]string, len(record))
	for k, v := range record {
		recordToWrite[k] = c.quoteField(orm.ToStr(v))
	}
	err = c.writer.Write(recordToWrite)
	return err
}

func (c *CSVDataFile) quoteField(field string) string {
	if strings.IndexAny(field, `\`) > 0 {
		str := strings.Replace(field, `\`, `\\`, -1)
		return str
	}
	return field
}

func (c *CSVDataFile) Flush() error {
	c.writer.Flush()
	return c.writer.Error()
}

func (c *CSVDataFile) Close() error {
	return c.fh.Close()
}

func (c *CSVDataFile) Remove() error {
	if !c.hasImportError {
		return os.Remove(c.File)
	} else {
		return fmt.Errorf("csvdatafile import error. keep csv file for debug.")
	}
}

func (c *CSVDataFile) LoadToMySQL(o orm.Ormer) error {
	sql := fmt.Sprintf(`load data infile '%s' into table %s FIELDS TERMINATED BY '%s' ENCLOSED BY '"' (%s)`, c.File, c.Table, string(c.FieldsTerminatedBy), strings.Join(c.Fields, ", "))
	//fmt.Println(sql)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		c.hasImportError = true
		return fmt.Errorf("datafile %s: %s", c.File, err.Error())
	}
	return nil
}
