package adoquery

import (
	"database/sql"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Driver int32

const (
	SQLServer Driver = 1 << (8 - 1 - iota)
	MySql
	Sqlite
	Postgresql
	Oracle
)

type Database struct {
	// 数据库类型
	Driver Driver
	// 连接字符串
	Path     string
	Port     string
	Username string
	Password string
	Dbname   string
	Config   string
	// 数据库实例
	db *gorm.DB
}

func (d *Database) Query(sql string, value ...any) (rows *sql.Rows, err error) {
	tx := d.db.Raw(sql, value...)
	if tx.Error != nil {
		err = tx.Error
	}
	rows, err = tx.Rows()
	return
}

func (d *Database) Connect() (err error) {
	var dsn string
	switch d.Driver {
	case SQLServer:
		dsn = "sqlserver://" + d.Username + ":" + d.Password + "@" + d.Path + ":" + d.Port + "?database=" + d.Dbname + "&encrypt=disable"
		d.db, err = newMssql(dsn)
	case MySql:
		dsn = d.Username + ":" + d.Password + "@tcp(" + d.Path + ":" + d.Port + ")/" + d.Dbname + "?" + d.Config
		d.db, err = newMysql(dsn)
	case Sqlite:
		dsn = filepath.Join(d.Path, d.Dbname+".db")
		d.db, err = newSqlite(dsn)
	case Postgresql:
		dsn = "host=" + d.Path + " user=" + d.Username + " password=" + d.Password + " dbname=" + d.Dbname + " port=" + d.Port + " " + d.Config
		d.db, err = newPgsql(dsn)
	case Oracle:
		dsn = "oracle://" + d.Username + ":" + d.Password + "@" + d.Path + ":" + d.Port + "/" + d.Dbname + "?" + d.Config
		d.db, err = newOracle(dsn)
	default:
		dsn = "sqlserver://" + d.Username + ":" + d.Password + "@" + d.Path + ":" + d.Port + "?database=" + d.Dbname + "&encrypt=disable"
		d.db, err = newMssql(dsn)
	}
	return
}

func (d *Database) Disconnect() {
	if d.db != nil {
		db, _ := d.db.DB()
		db.Close()
	}
}

func newMssql(dsn string) (db *gorm.DB, err error) {
	config := sqlserver.Config{
		DSN:               dsn,
		DefaultStringSize: 191,
	}
	db, err = gorm.Open(sqlserver.New(config), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})

	return
}

func newMysql(dsn string) (db *gorm.DB, err error) {
	config := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	db, err = gorm.Open(mysql.New(config), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	return
}

func newSqlite(dsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	return
}

func newPgsql(dsn string) (db *gorm.DB, err error) {
	config := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	}
	db, err = gorm.Open(postgres.New(config), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	return
}

func newOracle(dsn string) (db *gorm.DB, err error) {
	config := mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 191,
	}
	db, err = gorm.Open(mysql.New(config), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	return
}
