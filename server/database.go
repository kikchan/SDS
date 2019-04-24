package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

/*
	CREATES A NEW USER AND INITIALIZE THE PASSWORDS, CARDS AND NOTES
	TABLES WITH DEFAULT VALUES FOR THAT USER.
	THAT LAST PART IS DONE BY A TRIGGER ON THE MYSQL SERVER.

	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func createUser(user string, password string, hash string, salt string, data string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	var query = "INSERT INTO users(user, password, hash, salt, data) " +
		"VALUES ('" + user + "', '" + password + "', '" + hash + "', '" + salt + "', '" + data + "');"

	insert, err := db.Query(query)

	if err != nil {
		code = -2
		msg = "Error creating a new user: " + user
	} else {
		code = 1
		msg = "User created: " + user

		defer insert.Close()
	}

	return code, msg
}
