package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func writeLog(username string, msg string) {
	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		fmt.Println("Cannot connect to the database")
	}

	defer db.Close()

	var query = "INSERT INTO log(user, date, msg) VALUES ('" + username + "', NOW()" + ", \"" + strings.Replace(msg, "'", "", -1) + "\");"

	insert, err := db.Query(query)

	if err == nil {
		defer insert.Close()
	}
}
