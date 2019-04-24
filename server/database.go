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
	   -1: Error connecting to database
	   -2: Error executing query
*/
func createUser(user string, password string, hash string, salt string, data string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	var query = "INSERT INTO users(username, password, hash, salt, data) " +
		"VALUES ('" + user + "', '" + password + "', '" + hash + "', '" + salt + "', '" + data + "');"

	insert, err := db.Query(query)

	if err != nil {
		code = -2
		msg = err.Error()
	} else {
		code = 1
		msg = "User created: " + user

		defer insert.Close()
	}

	return code, msg
}

/*
	SEARCHES FOR AN EXISTING USER IN THE DATABASE AND RETURNS ALL THE COLUMNS.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
*/
func findUser(username string) (int, string, user) {
	var msg string
	var code int
	var user user

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	var query = "SELECT * FROM users WHERE username='" + username + "';"

	read, err := db.Query(query)
	if err != nil {
		code = -2
		msg = err.Error()
	}

	defer read.Close()

	if read.Next() {
		var a int
		var b, c, d, e, f string

		err = read.Scan(&a, &b, &c, &d, &e, &f)

		code = 1
		msg = "Found user: " + username
		user.ID = a
		user.Username = b
		user.Password = c
		user.Hash = decode64(d)
		user.Salt = decode64(e)
		user.Data = f
	} else {
		code = -3
		msg = "The user \"" + username + "\" doesn't exist"
	}

	return code, msg, user
}

/*
	DELETES AN EXISTING USER GIVEN ITS USERNAME.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
*/
func deleteUser(username string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(username)

	if code == 1 {
		var query = "DELETE FROM users WHERE username='" + username + "';"

		delete, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			code = 1
			msg = "User deleted: " + username

			defer delete.Close()
		}
	}

	return code, msg
}

/*
	UPDATE PASSWORDS' DATA FIELD WITH NEW INFORMATION.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
*/
func updatePassword(user string, data string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	var query = "UPDATE passwords SET data='" + data + "' WHERE user='" + user + "';"

	update, err := db.Query(query)
	if err != nil {
		code = -2
		msg = err.Error()
	} else {
		code = 1
		msg = "Passwords modified for user: " + user

		defer update.Close()
	}

	return code, msg
}

/*
	UPDATE CARDS' DATA FIELD WITH NEW INFORMATION.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
*/
func updateCard(user string, data string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	var query = "UPDATE cards SET data='" + data + "' WHERE user='" + user + "';"

	update, err := db.Query(query)
	if err != nil {
		code = -2
		msg = err.Error()
	} else {
		code = 1
		msg = "Cards modified for user: " + user

		defer update.Close()
	}

	return code, msg
}

/*
	UPDATE NOTES' DATA FIELD WITH NEW INFORMATION.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
*/
func updateNote(user string, data string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	var query = "UPDATE notes SET data='" + data + "' WHERE user='" + user + "';"

	update, err := db.Query(query)
	if err != nil {
		code = -2
		msg = err.Error()
	} else {
		code = 1
		msg = "Notes modified for user: " + user

		defer update.Close()
	}

	return code, msg
}
