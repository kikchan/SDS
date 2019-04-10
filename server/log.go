package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func writeLog(msg string) {
	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		fmt.Println("Cannot connect to the database")
	}

	defer db.Close()

	//fmt.Println(strings.Replace(msg, "'", "", -1))

	var query = "INSERT INTO log VALUES (NOW()" + "', '" + strings.Replace(msg, "'", "", -1) + ");"

	fmt.Println(query)

	insert, err := db.Query(query)

	if err == nil {
		defer insert.Close()
	}
}
