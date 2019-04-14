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
func createPassword(username string, pass string, user string) (int, string) {
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
		var query = "INSERT INTO passwords(username, pass, date, user) VALUES ('" + username + "', '" + pass + "', NOW(), '" + user + "');"
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
	READ a note by its ID
	Returns:
		1: OK
	   -1: Invalid note
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

	code, msg = findUser(user)

	if code == 1 {
		var query = "SELECT * FROM notes WHERE user='" + user + "' AND id=" + strconv.Itoa(id) + ";"
		writeLog(user, "findNoteByID", query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		}

		defer read.Close()

		if read.Next() {
			var a, b, c, d string

			err = read.Scan(&a, &b, &c, &d)

			code = 1
			msg = a + " " + b + " " + c + " " + d
		} else {
			code = -1
			msg = "Invalid note"
		}
	}

	writeLog(user, "findNoteByID response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	READ all notes
	Returns:
		1: OK
	   -1: The user doesn't have any notes
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func getUserPasswords(user string) (int, string) {
	var msg string
	var code int
	var notes []string

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg = findUser(user)

	if code == 1 {
		var query = "SELECT * FROM notes WHERE user='" + user + "';"

		writeLog(user, "getUserNotes", query)

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
			notes = append(notes, "["+a+" "+b+" "+c+" "+d+"]")
		}

		if len(notes) != 0 {
			code = 1
			msg = ""

			for i := 0; i < len(notes); i++ {
				msg += notes[i]
			}
		} else {
			code = -1
			msg = "The user has no notes"
		}
	}

	writeLog(user, "getUserNotes response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

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

	code, msg = findUser(user)

	if code == 1 {
		code, msg = findNoteByID(user, id)

		if code == 1 {
			var query = "UPDATE notes SET text=\"" + text + "\" WHERE user='" + user + "' AND id=" + strconv.Itoa(id) + ";"
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

	code, msg = findUser(user)

	if code == 1 {
		code, msg = findNoteByID(user, id)

		if code == 1 {
			var query = "DELETE FROM notes WHERE user='" + user + "' AND id=" + strconv.Itoa(id) + ";"
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
