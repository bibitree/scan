package model

import (
	"database/sql"
	"fmt"
)

type mysqlPool struct {
	*sql.DB
	Namespace string
}

var MysqlPool *mysqlPool

func InitMysql(c *Config) error {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Address, c.Database)

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
