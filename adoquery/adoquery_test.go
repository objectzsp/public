package adoquery

import (
	"fmt"
	"testing"
)

func TestADOQuery(t *testing.T) {
	connection := Database{
		Driver:   SQLServer,
		Path:     "zs2019.gdfzjy.com",
		Port:     "9304",
		Username: "sa",
		Password: "Sl81262299",
		Dbname:   "zs2020_backup",
	}
	ado := New(connection)
	ado.SQL("select top 1 * from F_BH_zdbpbh where bpbhid > ?", "1", 2, '3')
	ado.Open()
	fmt.Println(ado.Sql)
	fmt.Println(string(ado.Json))
	ado.Close()
}
