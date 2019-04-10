package main

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

/*
	CREATE
	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func createCard(pan string, ccv string, month int, year int, owner string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	var query = "INSERT INTO cards(pan, ccv, expiry, `owner`) VALUES ('" + pan + "', '" + ccv + "', '" + strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/00', '" + owner + "');"
	writeLog(owner, "[Query]: "+query)

	insert, err := db.Query(query)
	if err != nil {
		code = -2
		msg = "Invalid card"
	} else {
		code = 1
		msg = "Added new card(" + pan + ") for user: " + owner

		defer insert.Close()
	}

	writeLog(owner, "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	READ
	Returns:
		1: OK
	   -1: Card doesn't exist
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func findCardByPAN(owner string, pan string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	var query = "SELECT * FROM cards WHERE owner='" + owner + "' AND pan='" + pan + "';"
	writeLog(owner, "[Query]: "+query)

	read, err := db.Query(query)
	if err != nil {
		code = -2
		msg = err.Error()
	}

	defer read.Close()

	if read.Next() {
		var a, b, c, d, e string

		err = read.Scan(&a, &b, &c, &d, &e)

		code = 1
		msg = a + " " + b + " " + c + " " + d + " " + e + " "
	} else {
		code = -1
		msg = "Invalid card"
	}

	writeLog(owner, "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	READ
	Returns:
		1: OK
	   -1: Card doesn't exist
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func findCardByID(owner string, id int) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	var query = "SELECT * FROM cards WHERE owner='" + owner + "' AND id=" + strconv.Itoa(id) + ";"
	writeLog(owner, "[Query]: "+query)

	read, err := db.Query(query)
	if err != nil {
		code = -2
		msg = err.Error()
	}

	defer read.Close()

	if read.Next() {
		var a, b, c, d, e string

		err = read.Scan(&a, &b, &c, &d, &e)

		code = 1
		msg = a + " " + b + " " + c + " " + d + " " + e + " "
	} else {
		code = -1
		msg = "Invalid card"
	}

	writeLog(owner, "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	UPDATE
	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func updateCard(username string, password string, email string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	update, err := db.Query("UPDATE users SET password='" + password + "', email='" + email + "' WHERE username='" + username + "';")
	if err != nil {
		code = -2
		msg = err.Error()
	} else {
		code = 1
		msg = "User modified: " + username
	}

	defer update.Close()

	return code, msg
}

/*
	DELETE
	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func deleteCard(username string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	delete, err := db.Query("DELETE FROM users WHERE username='" + username + "';")
	if err != nil {
		code = -2
		msg = err.Error()
	} else {
		code = 1
		msg = "User deleted: " + username
	}

	defer delete.Close()

	return code, msg
}
