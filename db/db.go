package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DatabaseConnection() (db *sql.DB) {
	db, err := sql.Open("mysql", "root@tcp(db:3306)/go?parseTime=true")

	if err != nil {
		panic(err.Error())
	}

	return db
}
