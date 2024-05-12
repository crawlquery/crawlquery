package testutil

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func CreateTestMysqlDB() *sql.DB {
	db, err := sql.Open("mysql", "root:cqdb@tcp(localhost:3306)/cqdb_test?parseTime=true")

	if err != nil {
		panic(err)
	}

	return db
}
