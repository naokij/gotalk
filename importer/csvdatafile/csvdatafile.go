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
		recordToWrite[k] = orm.ToStr(v)
	}
	err = c.writer.Write(recordToWrite)
	return err
}

func (c *CSVDataFile) Flush() error {
	c.writer.Flush()
	return c.writer.Error()
}

func (c *CSVDataFile) Close() error {
	return c.fh.Close()
}

func (c *CSVDataFile) Remove() error {
	return os.Remove(c.File)
}

func (c *CSVDataFile) LoadToMySQL(o orm.Ormer) error {
	sql := fmt.Sprintf(`load data infile '%s' into table %s FIELDS TERMINATED BY '%s' ENCLOSED BY '"' (%s)`, c.File, c.Table, string(c.FieldsTerminatedBy), strings.Join(c.Fields, ", "))
	//fmt.Println(sql)
	_, err := o.Raw(sql).Exec()
	return err
}
