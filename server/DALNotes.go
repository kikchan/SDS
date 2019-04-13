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
	READ a card by its ID
	Returns:
		1: OK
	   -1: Invalid card
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func findNoteByID(user string, id int) (int, string) {
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
	READ all cards
	Returns:
		1: OK
	   -1: The user doesn't have any cards
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func getUserNotes(owner string) (int, string) {
	var msg string
	var code int
	var cards []string

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg = findUser(owner)

	if code == 1 {
		var query = "SELECT * FROM cards WHERE owner='" + owner + "';"
		writeLog(owner, "getUserCards", query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		}

		defer read.Close()

		for read.Next() {
			var a, b, c, d, e string

			err = read.Scan(&a, &b, &c, &d, &e)

			code = 1
			cards = append(cards, "["+a+" "+b+" "+c+" "+d+" "+e+"]")
		}

		if len(cards) != 0 {
			code = 1

			for i := 0; i < len(cards); i++ {
				msg += cards[i]
			}
		} else {
			code = -1
			msg = "The user has no cards"
		}
	}

	writeLog(owner, "getUserCards response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	UPDATE
	Returns:
		1: OK
	   -1: Invalid card
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func updateNote(pan string, ccv string, month int, year int, owner string, oldPAN string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg = findUser(owner)

	if code == 1 {
		code, msg = findCardByPAN(owner, pan)

		if code == 1 {
			var query = "UPDATE cards SET pan='" + pan + "', ccv='" + ccv + "', expiry='" + strconv.Itoa(year) + "/" +
				strconv.Itoa(month) + "/00' WHERE owner='" + owner + "' AND pan='" + oldPAN + "';"
			writeLog(owner, "updateCard", query)

			update, err := db.Query(query)
			if err != nil {
				code = -2
				msg = err.Error()
			} else {
				code = 1
				msg = "Card modified: " + pan

				defer update.Close()
			}
		} else {
			code = -1
			msg = "Invalid card: " + pan
		}
	}

	writeLog(owner, "updateCard response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}

/*
	DELETE
	Returns:
		1: OK
	   -1: Invalid card
	   -2: Error executing query
	   -3: Error connecting to DB
*/
func deleteNote(pan string, owner string) (int, string) {
	var msg string
	var code int

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -3
		msg = err.Error()
	}

	defer db.Close()

	code, msg = findUser(owner)

	if code == 1 {
		code, msg = findCardByPAN(owner, pan)

		if code == 1 {
			var query = "DELETE FROM cards WHERE owner='" + owner + "' AND pan='" + pan + "';"
			writeLog(owner, "deleteCard", query)

			if code == 1 {
				delete, err := db.Query(query)
				if err != nil {
					code = -2
					msg = err.Error()
				} else {
					code = 1
					msg = "Card deleted: " + pan

					defer delete.Close()
				}
			} else {
				code = -1
				msg = "Invalid card"
			}
		}
	}

	writeLog(owner, "deleteCard response", "[Result]: code: "+strconv.Itoa(code)+" ## msg: "+msg)

	return code, msg
}
