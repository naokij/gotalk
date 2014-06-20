package converters

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	//"time"
)

type TableRows struct {
	Rows int64
}

func NumOfRows(o orm.Ormer, sql string) (int64, error) {
	type tableRows struct {
		Rows int64
	}
	var rows tableRows
	if err := o.Raw(sql).QueryRow(&rows); err != nil {
		return 0, err
	}
	return rows.Rows, nil
}

func RunPreImportMySQLSettings(o orm.Ormer) {
	const sql = `/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;`
	var cmds = strings.Split(sql, "\n")
	for _, cmd := range cmds {
		if _, err := o.Raw(cmd).Exec(); err != nil {
			fmt.Println("PreImportSettings: ", err)
		}
	}
}

func cutString(str string, length int) string {
	chars := []rune(str)
	if len(chars) <= length {
		return str
	}
	return string(chars[0:length])
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Map2InsertSql(o orm.Ormer, table string, data map[string]interface{}) error {
	var values []interface{}
	var keys []string
	var valuesPlaceHolders []string
	for k, v := range data {
		valuesPlaceHolders = append(valuesPlaceHolders, "? ")
		keys = append(keys, fmt.Sprintf("`%s`", k))
		values = append(values, v) //
		// switch v.(type) {
		// case bool:
		// 	values = append(values, fmt.Sprintf("%d", Btoi(v.(bool))))
		// case int:
		// 	values = append(values, fmt.Sprint(v.(int)))
		// case int32:
		// 	values = append(values, fmt.Sprint(v.(int32)))
		// case int64:
		// 	values = append(values, fmt.Sprint(v.(int64)))
		// case float64:
		// 	values = append(values, fmt.Sprint(v.(float64)))
		// case string:
		// 	values = append(values, v.(string))
		// case time.Time:
		// 	t := v.(time.Time)
		// 	formattedTime := t.Format("2006-01-02 15:04:05")
		// 	values = append(values, formattedTime)
		// case nil:
		// 	values = append(values, "NULL")
		// }
	}
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES(%s)", table, strings.Join(keys, ","), strings.Join(valuesPlaceHolders, ","))
	_, err := o.Raw(sql, values...).Exec()
	if err != nil {
		return err
	}
	return nil
}
