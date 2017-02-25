package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhangmingkai4315/weichat-golang-backend/config"
)

// NewDatabase will create a new connection and return it for save database
// It will return a sql.DB object and a error when connection failure.
func NewDatabase() (*sql.DB, error) {
	databaseURL := fmt.Sprintf("%s:%s@/%s",
		config.ConfigObj.Database.User,
		config.ConfigObj.Database.Password,
		config.ConfigObj.Database.DBName)
	log.Println(config.ConfigObj)
	database, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Panicf("Create Database Connection with Error %s", err.Error())
		return nil, err
	}
	err = database.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
		return nil, err
	}

	return database, nil
}
