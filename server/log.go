package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

/*
	WRITES THE INPUT OF THE REQUEST AND THE OUTPUT OF THE RESPONSE FROM A FUN
*/
func writeLog(method string, username string, id int, msg string) {
	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		consoleTimeStamp()
		fmt.Println(err.Error())
	}

	defer db.Close()

	var query = "INSERT INTO log(date, method, user, correlation, text) " +
		"VALUES(NOW(), '" + method + "', '" + username + "', '" + strconv.Itoa(id) + "', \"" + strings.Replace(msg, "'", "", -1) + "\");"

	insert, err := db.Query(query)

	if err == nil {
		defer insert.Close()
	}
}
