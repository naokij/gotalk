package converters

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
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

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
