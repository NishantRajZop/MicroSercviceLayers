package dataSource

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB(filepath string) *sql.DB {
	var err error
	db, err := sql.Open("mysql", filepath)
	if err != nil {
		panic(err)
	}
	return db
}
