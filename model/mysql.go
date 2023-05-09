package model

import (
	"database/sql"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type mysqlPool struct {
	*sql.DB
	Namespace string
}

var MysqlPool *mysqlPool

func InitMysql(c *Config) error {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		c.Username, c.Password, c.Address, c.Database)

	// dataSourceName := "root:a2683570@tcp(127.0.0.1:3306)/chainfindata?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(c.MaxIdle)
	db.SetMaxOpenConns(c.MaxActive)

	MysqlPool = &mysqlPool{
		Namespace: c.Namespace,
		DB:        db,
	}

	err = MysqlPool.Ping()
	if err != nil {
		return err
	}

	return nil
}

func DisposeMysql() {
	if MysqlPool != nil {
		MysqlPool.Close()
	}
}
