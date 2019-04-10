package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

/*
	CREATE
	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func createUser(username string, password string, name string, surname string, email string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	insert, err := db.Query("INSERT INTO users VALUES ('" + username + "', '" + password + "', '" + name + "', '" + surname + "', '" + email + "');")
	if err != nil {
		code = -2
		msg = "Invalid username"
	} else {
		code = 1
		msg = "User created: " + username + ", " + name + " " + surname

		defer insert.Close()
	}

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
func findUser(username string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	read, err := db.Query("SELECT * FROM users WHERE username='" + username + "';")
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

	return code, msg
}

/*
	UPDATE
	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func updateUser(username string, password string, email string) (int, string) {
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
		update, err := db.Query("UPDATE users SET password='" + password + "', email='" + email + "' WHERE username='" + username + "';")
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

	return code, msg
}

/*
	DELETE
	Returns:
		1: OK
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func deleteUser(username string) (int, string) {
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
		delete, err := db.Query("DELETE FROM users WHERE username='" + username + "';")
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

	return code, msg
}
