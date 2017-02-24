package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhangmingkai4315/weichat-golang-backend/config"
)

func NewDatabase() (*sql.DB, error) {
	databaseUrl := fmt.Sprintf("%s:%s@/%s",
		config.ConfigObj.Database.User,
		config.ConfigObj.Database.Password,
		config.ConfigObj.Database.DBName)
	database, err := sql.Open("mysql", databaseUrl)
	if err != nil {
		//log.Panicf("Create Database Connection with Error %s", err.Error())
		return nil, err
	}
	err = database.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		//log.Fatalf("Error on opening database connection: %s", err.Error())
		return nil, err
	}

	return database, nil
}
