package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DatabaseConnection() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbName := "go"

	db, err := sql.Open(dbDriver, dbUser+"@/"+dbName+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}

	return db
}
