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
func createUser(username string, password string, hash string, salt string, name string, surname string, email string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	var query = "INSERT INTO users VALUES ('" + username + "', '" + password + "', '" + hash + "', '" + salt + "', '" + name + "', '" + surname + "', '" + email + "');"
	writeLog(username, "createUser", query)

	insert, err := db.Query(query)

	if err != nil {
		code = -2
		msg = "Invalid username"
	} else {
		code = 1
		msg = "User created: " + username + ", " + name + " " + surname

		defer insert.Close()
	}

	writeLog(username, "createUser response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

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
func findUser(username string) (int, string, user) {
	var msg string
	var code int
	var user user

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	var query = "SELECT * FROM users WHERE username='" + username + "';"
	writeLog(username, "findUser", query)

	read, err := db.Query(query)
	if err != nil {
		code = -2
		msg = err.Error()
	}

	defer read.Close()

	if read.Next() {
		var a, b, c, d, e, f, g string

		err = read.Scan(&a, &b, &c, &d, &e, &f, &g)

		code = 1
		msg = a + " " + b + " " + c + " " + d + " " + e + " " + " " + f + " " + g
		user.Username = a
		user.Password = b
		user.Hash = decode64(c)
		user.Salt = decode64(d)
		user.Name = e
		user.Surname = f
		user.Email = g
	} else {
		code = -1
		msg = "Invalid username"
	}

	writeLog(username, "findUser response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg, user
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

	code, msg, _ = findUser(username)

	if code == 1 {
		var query = "UPDATE users SET password='" + password + "', email='" + email + "' WHERE username='" + username + "';"
		writeLog(username, "updateUser", query)

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

	writeLog(username, "updateUser response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

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

	code, msg, _ = findUser(username)

	if code == 1 {
		var query = "DELETE FROM users WHERE username='" + username + "';"
		writeLog(username, "deleteUser", query)

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

	writeLog(username, "deleteUser response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}
