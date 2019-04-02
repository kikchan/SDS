package cad

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Export struct {
}

// DB asdasdas
func (c Export) DB() {
	db, err := sql.Open("mysql", "sds:sds@tcp(185.207.145.237:3306)/sds")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//insert, err := db.Query("INSERT INTO sds.users (username, password, name, surname) VALUES (sds, sds, SDS, SDS);")
	insert, err := db.Query("INSERT INTO users VALUES ('sds', 'sds', 'SDS', 'SDS');")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}
