package utils

import "database/sql"

var (
	Db  *sql.DB
	err error
)

func init() {
	Db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/bookstore")
	if err != nil {
		panic(err.Error())
	}
}
