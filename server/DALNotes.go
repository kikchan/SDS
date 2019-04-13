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
func createNote(text string, user string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg = findUser(user)

	if code == 1 {
		var query = "INSERT INTO notes(text, modified, user) VALUES (\"" + text + "\", NOW(), '" + user + "');"
		writeLog(user, "createNote", query)

		insert, err := db.Query(query)

		if err != nil {
			code = -2
			msg = "Error executing the query"
		} else {
			code = 1
			msg = "Note created for user: " + user

			defer insert.Close()
		}
	}

	writeLog(user, "createNote response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	READ
	Returns:
		1: OK
	   -1: User doesn't exist
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func findNote(username string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	var query = "SELECT * FROM users WHERE username='" + username + "';"
	//writeLog(username, "[Function]: findUser ## [Query]: "+query)

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
		msg = "Invalid username"
	}

	//writeLog(username, "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	UPDATE
	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func updateNote(username string, password string, email string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg = findUser(username)

	if code == 1 {
		var query = "UPDATE users SET password='" + password + "', email='" + email + "' WHERE username='" + username + "';"
		//writeLog(username, "[Function]: updateUser ## [Query]: "+query)

		update, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			code = 1
			msg = "User modified: " + username

			defer update.Close()
		}
	} else {
		code = -2
		msg = "Invalid username"
	}

	//writeLog(username, "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	DELETE
	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func deleteNote(username string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg = findUser(username)

	if code == 1 {
		var query = "DELETE FROM users WHERE username='" + username + "';"
		//writeLog(username, "[Function]: deleteUser ## [Query]: "+query)

		delete, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			code = 1
			msg = "User deleted: " + username

			defer delete.Close()
		}
	} else {
		code = -2
		msg = "User \"" + username + "\" doesn't exist"
	}

	//writeLog(username, "[Function]: deleteUser ## [Query]: "+query)

	return code, msg
}
