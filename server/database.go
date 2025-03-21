package main

import (
	"database/sql"
	"math/rand"

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
func createUser(user string, pubKey string, hash string, salt string, data string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	var query = "INSERT INTO users(username, publicKey, hash, salt, data) " +
		"VALUES ('" + user + "', '" + pubKey + "', '" + hash + "', '" + salt + "', '" + data + "');"

	writeLog("createUser entry", user, correlation, query)

	insert, err := db.Query(query)

	if err != nil {
		code = -2
		msg = err.Error()
	} else {
		code = 1
		msg = "User created: " + user

		defer insert.Close()
	}

	writeLog("createUser out", user, correlation, msg)

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
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	var query = "SELECT * FROM users WHERE username='" + username + "';"

	writeLog("findUser entry", username, correlation, query)

	read, err := db.Query(query)
	if err != nil {
		code = -2
		msg = err.Error()
	}

	defer read.Close()

	if read.Next() {
		var a int
		var b, d, e, f, g string

		err = read.Scan(&a, &b, &d, &e, &f, &g)

		code = 1
		msg = "Found user: " + username
		user.ID = a
		user.Username = b
		user.PubKey = d
		user.Hash = decode64(e)
		user.Salt = decode64(f)
		user.Data = g
	} else {
		code = -3
		msg = "The user \"" + username + "\" doesn't exist"
	}

	writeLog("findUser out", username, correlation, msg)

	return code, msg, user
}

/*
	UPDATE DATA USER' DATA FIELD WITH NEW INFORMATION.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
*/
func updateDataUser(username string, data string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(username)

	if code == 1 {
		var query = "UPDATE users SET data='" + data + "' WHERE username='" + username + "';"

		writeLog("updateDataUser entry", username, correlation, query)

		update, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			msg = "Data user modified for username: " + username

			defer update.Close()
		}
	}

	writeLog("updateDataUser out", username, correlation, msg)

	return code, msg
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
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(username)

	if code == 1 {
		var query = "DELETE FROM users WHERE username='" + username + "';"

		writeLog("deleteUser entry", username, correlation, query)

		delete, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			msg = "User deleted: " + username

			defer delete.Close()
		}
	}

	writeLog("deleteUser out", username, correlation, msg)

	return code, msg
}

/*
	UPDATE PASSWORDS' DATA FIELD WITH NEW INFORMATION.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
*/
func updatePassword(user string, data string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "UPDATE passwords SET data='" + data + "' WHERE user='" + user + "';"

		writeLog("updatePassword entry", user, correlation, query)

		update, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			msg = "Passwords modified for user: " + user

			defer update.Close()
		}
	}

	writeLog("updatePassword out", user, correlation, msg)

	return code, msg
}

/*
	GET PASSWORDS FOR A GIVEN USER

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
	   -4: No passwords were found
*/
func getUserPasswords(user string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "SELECT data FROM passwords WHERE user='" + user + "';"

		writeLog("getUserPasswords entry", user, correlation, query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			defer read.Close()

			code = -4
			msg = "No passwords were found"

			if read.Next() {
				var a string

				err = read.Scan(&a)

				code = 1
				msg = a
			}
		}
	}

	writeLog("getUserPasswords out", user, correlation, msg)

	return code, msg
}

/*
	UPDATE CARDS' DATA FIELD WITH NEW INFORMATION.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
*/
func updateCard(user string, data string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "UPDATE cards SET data='" + data + "' WHERE user='" + user + "';"

		writeLog("updateCard entry", user, correlation, query)

		update, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			msg = "Cards modified for user: " + user

			defer update.Close()
		}
	}

	writeLog("updateCard out", user, correlation, msg)

	return code, msg
}

/*
	GET CARDS FOR A GIVEN USER

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
	   -5: No cards were found
*/
func getUserCards(user string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "SELECT data FROM cards WHERE user='" + user + "';"

		writeLog("getUserCards entry", user, correlation, query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			defer read.Close()

			code = -4
			msg = "No cards were found"

			if read.Next() {
				var a string

				err = read.Scan(&a)

				code = 1
				msg = a
			}
		}
	}

	writeLog("getUserCards out", user, correlation, msg)

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
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "UPDATE notes SET data='" + data + "' WHERE user='" + user + "';"

		writeLog("updateNote entry", user, correlation, query)

		update, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			msg = "Notes modified for user: " + user

			defer update.Close()
		}
	}

	writeLog("updateNote out", user, correlation, msg)

	return code, msg
}

/*
	GET NOTES FOR A GIVEN USER

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
	   -6: No notes were found
*/
func getUserNotes(user string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "SELECT data FROM notes WHERE user='" + user + "';"

		writeLog("getUserNotes entry", user, correlation, query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			defer read.Close()

			code = -4
			msg = "No notes were found"

			if read.Next() {
				var a string

				err = read.Scan(&a)

				code = 1
				msg = a
			}
		}
	}

	writeLog("getUserNotes out", user, correlation, msg)

	return code, msg
}

/*
	RETREIVES ALL THE USERS EXCEPT THE GIVEN ONE.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
	   -8: No users available to share
*/
func getAllUsersExceptTheGivenOne(username string) (int, string, []user) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)
	var users []user

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(username)

	if code == 1 {
		var query = "SELECT username, publicKey FROM users WHERE username<>'" + username + "';"

		writeLog("getAllUsersExceptTheGivenOne entry", username, correlation, query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			defer read.Close()

			code = -8
			msg = "No users available to share"

			for read.Next() {
				var a, b string

				err = read.Scan(&a, &b)

				var us user
				us.Username = a
				us.PubKey = b

				users = append(users, us)

				code = 1
				msg = "Other users to share with were found"
			}
		}
	}

	writeLog("getAllUsersExceptTheGivenOne out", username, correlation, msg)

	return code, msg, users
}

/*
	SHARES AN ENCRYPTED FIELD WITH A NUMBER OF USERS.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
*/
func shareField(user, typeF, fieldID, data, userTarget, userKey string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "INSERT INTO shares(user, type, fieldId, data, user_target, user_key) " +
			"VALUES ('" + user + "', '" + typeF + "'," + fieldID + ", '" + data + "', '" + userTarget + "', '" + userKey + "');"

		writeLog("shareField entry", user, correlation, query)

		insert, err := db.Query(query)

		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			defer insert.Close()

			code = 1
			msg = "Field successfully shared"
		}
	}

	writeLog("shareField out", user, correlation, msg)

	return code, msg
}

/***********************************************************************************
/***************************	PARA REVISAR	************************************
/***********************************************************************************

/*
	RETREIVES A SHARED FIELD FOR A GIVEN USER.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
	   -9: No fields were found
*/
func getSharedFieldsForUser(user string, typeF string) (int, string, []field) {
	var msg string
	var code int
	var fields []field
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "SELECT `data`, user_key FROM shares WHERE user_target='" + user + "' AND type='" + typeF + "';"

		writeLog("getSharedFieldsForUser entry", user, correlation, query)

		read, err := db.Query(query)

		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			defer read.Close()

			code = -9
			msg = "No fields were found"

			for read.Next() {
				var a, b string
				var field field

				err = read.Scan(&a, &b)

				code = 1
				msg = "Retrieved fields for user: \"" + user + "\""

				field.Data = a
				field.UserKey = b

				fields = append(fields, field)
			}
		}
	}

	writeLog("getSharedFieldsForUser out", user, correlation, msg)

	return code, msg, fields
}

/*
	RETREIVES A SHARED FIELD.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
	   -7: The field doesn't exist
*/
func sharedFieldExists(user, typeF, fieldID string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		var query = "SELECT * FROM shares WHERE user='" + user + "' AND type='" + typeF +
			"' AND fieldId=" + fieldID + ";"

		writeLog("sharedFieldExists entry", user, correlation, query)

		read, err := db.Query(query)
		if err != nil {
			code = -2
			msg = err.Error()
		} else {
			defer read.Close()

			if read.Next() {
				code = 1
				msg = "The requested field is shared and it's going to be re-shared again"
			} else {
				code = -7
				msg = "The field doesn't exist"
			}
		}
	}

	writeLog("sharedFieldExists out", user, correlation, msg)

	return code, msg
}

/*
	UPDATES A SHARED FIELD.

	Returns:
		1: OK
	   -1: Error connecting to database
	   -2: Error executing query
	   -3: The user doesn't exist
	   -7: The field doesn't exist
*/
func updateSharedField(user, typeF, fieldID, data string) (int, string) {
	var msg string
	var code int
	var correlation = rand.Intn(10000)

	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	if err != nil {
		code = -1
		msg = err.Error()
	}

	defer db.Close()

	code, msg, _ = findUser(user)

	if code == 1 {
		code, msg = sharedFieldExists(user, typeF, fieldID)

		if code == 1 {
			var query = "UPDATE shares SET data='" + data + "' WHERE user='" + user + "' AND type='" + typeF + "' AND fieldId=" + fieldID + ";"

			writeLog("updateSharedField entry", user, correlation, query)

			delete, err := db.Query(query)

			if err != nil {
				code = -2
				msg = err.Error()
			} else {
				defer delete.Close()

				code = 1
				msg = "Field successfully updated"
			}
		}
	}

	writeLog("updateSharedField out", user, correlation, msg)

	return code, msg
}
