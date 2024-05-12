package main

import (
	"database/sql"
	"fmt"
)

func main() {

	db, err := sql.Open("mysql", "root:cqdb@tcp(localhsot:3306)/cqdb")
	defer db.Close()

	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return
	}

	cqr := jobRepo.NewMysqlRepository(db)

}
