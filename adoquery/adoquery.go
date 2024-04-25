package adoquery

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type ADOQuery struct {
	// 数据连接对象
	Database *Database

	// 查询Sql
	Sql string

	// 响应结果
	Json         []byte
	RowsAffected int64
	Error        error
}

func (ado *ADOQuery) Open() {
	ado.AddError(ado.Database.Connect())
	if rows, err := ado.Database.Query(ado.Sql); err == nil {
		ado.RowsAffected = 0
		ado.ScanRows(rows)
		ado.AddError(rows.Close())
	}
}

func (ado *ADOQuery) ScanRows(rows *sql.Rows) {
	var (
		res         = make([]map[string]interface{}, 0)
		colTypes, _ = rows.ColumnTypes()
		value       = make([]interface{}, len(colTypes))
		parma       = make([]interface{}, len(colTypes))
	)

	for i, colType := range colTypes {
		value[i] = reflect.New(colType.ScanType())
		parma[i] = reflect.ValueOf(&value[i]).Interface()
	}

	for rows.Next() {
		ado.RowsAffected++
		rows.Scan(parma...)
		record := make(map[string]interface{})
		for i, colType := range colTypes {
			if value[i] == nil {
				record[colType.Name()] = ""
			} else {
				switch value[i].(type) {
				case []uint8:
					f, _ := strconv.ParseFloat(string(value[i].([]byte)), 64)
					record[colType.Name()] = f
				default:
					record[colType.Name()] = value[i]
				}
			}
		}
		res = append(res, record)
	}
	ado.Json, _ = json.Marshal(res)
}

func (ado *ADOQuery) SQL(sql string, vars ...interface{}) {

	var (
		idx       int
		sqlBuffer = bytes.NewBufferString("")
	)

	for _, b := range []byte(sql) {
		if b == '?' && len(vars) > idx {
			switch value := reflect.ValueOf(vars[idx]); value.Kind() {
			case reflect.String:
				sqlBuffer.WriteByte('\'')
				sqlBuffer.WriteString(value.String())
				sqlBuffer.WriteByte('\'')
			default:
				sqlBuffer.Write([]byte(fmt.Sprintf("%v", vars[idx])))
			}
		} else {
			sqlBuffer.WriteByte(b)
		}
	}
	// var (
	// 	idx int
	// 	sb  = bytes.NewBufferString("")
	// )
	// for _, v := range []byte(sql) {
	// 	if v == '?' && len(vars) > idx {
	// 		sb.Write([]byte(fmt.Sprintf("%v", vars[idx])))
	// 		idx++
	// 	} else {
	// 		sb.WriteByte(v)
	// 	}
	// }
	ado.Sql = sqlBuffer.String()
}

func (ado *ADOQuery) Close() {
	ado.Json = nil
	ado.Error = nil
	ado.RowsAffected = 0
	ado.Database.Disconnect()
}

func (ado *ADOQuery) AddError(err error) error {
	if err != nil {
		if ado.Error == nil {
			ado.Error = err
		} else {
			ado.Error = fmt.Errorf("%v; %w", ado.Error, err)
		}
	}
	return ado.Error
}

func (ADOQuery) GetJsonStr() string {
	return ""
}

func (ADOQuery) GetString(key string) (value string) {
	return ""
}

func New(conn Database) *ADOQuery {
	return &ADOQuery{
		Database: &conn,
	}
}
