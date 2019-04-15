package main

import (
	"database/sql"
	"fmt"
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
func createPassword(username string, pass string, user string, site string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "INSERT INTO passwords(username, pass, date, user, site) VALUES ('" + username + "', '" + pass + "', NOW(), '" + user + "', '" + site + "');"
		writeLog(user, "createPassword", query)

		insert, err := db.Query(query)

		if err != nil {
			code = -2
			msg = "Error executing the query"
		} else {
			code = 1
			msg = "Password added for user: " + user

			defer insert.Close()
		}
	}

	writeLog(user, "createPassword response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	READ a password given its ID
	Returns:
		1: OK
	   -1: Password not found
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func findPasswordByID(user string, id int) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "SELECT * FROM passwords WHERE user='" + user + "' AND id=" + strconv.Itoa(id) + ";"
		writeLog(user, "findPasswordByID", query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()

			fmt.Println(err.Error())
		}

		defer read.Close()

		if read.Next() {
			var a, b, c, d, e, f string

			err = read.Scan(&a, &b, &c, &d, &e, &f)

			code = 1
			msg = a + " " + b + " " + c + " " + d + " " + e + " " + f
		} else {
			code = -1
			msg = "The requested password was not found"
		}
	}

	writeLog(user, "findPasswordByID response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	READ all passwords
	Returns:
		1: OK
	   -1: The user doesn't have any passwords
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func getUserPasswords(user string) (int, string) {
	var msg string
	var code int
	var passwords []string

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "SELECT * FROM passwords WHERE user='" + user + "' ORDER BY site asc;"

		writeLog(user, "getUserPasswords", query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		}

		defer read.Close()

		for read.Next() {
			var a, b, c, d string

			err = read.Scan(&a, &b, &c, &d)

			code = 1
			passwords = append(passwords, "["+a+" "+b+" "+c+" "+d+"]")
		}

		if len(passwords) != 0 {
			code = 1
			msg = ""

			for i := 0; i < len(passwords); i++ {
				msg += passwords[i]
			}
		} else {
			code = -1
			msg = "The user has no passwords"
		}
	}

	writeLog(user, "getUserPasswords response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	UPDATE
	Returns:
		1: OK
	   -1: Invalid note
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func updatePassword(id int, text string, user string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		code, msg = findNoteByID(user, id)

		if code == 1 {
			var query = "UPDATE passwords SET text=\"" + text + "\" WHERE user='" + user + "' AND id=" + strconv.Itoa(id) + ";"
			writeLog(user, "updateNote", query)

			update, err := db.Query(query)
			if err != nil {
				code = -2
				msg = err.Error()
			} else {
				code = 1
				msg = "Note modified: " + strconv.Itoa(id)

				defer update.Close()
			}
		} else {
			code = -1
			msg = "Invalid note: " + strconv.Itoa(id)
		}
	}

	writeLog(user, "updateNote response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	DELETE
	Returns:
		1: OK
	   -1: Invalid note
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func deletePassword(id int, user string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		code, msg = findNoteByID(user, id)

		if code == 1 {
			var query = "DELETE FROM passwords WHERE user='" + user + "' AND id=" + strconv.Itoa(id) + ";"
			writeLog(user, "deleteNote", query)

			if code == 1 {
				delete, err := db.Query(query)
				if err != nil {
					code = -2
					msg = err.Error()
				} else {
					code = 1
					msg = "Note deleted: " + strconv.Itoa(id)

					defer delete.Close()
				}
			} else {
				code = -1
				msg = "Invalid note"
			}
		}
	}

	writeLog(user, "deleteNote response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}
