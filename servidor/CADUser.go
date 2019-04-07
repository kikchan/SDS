package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//CREATE
func createUser(username string, password string, name string, surname string, email string) string {
	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	insert, err := db.Query("INSERT INTO users VALUES ('" + username + "', '" + password + "', '" + name + "', '" + surname + "', '" + email + "');")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	return "Usuario creado: " + username + ", " + name + " " + surname
}

//READ
func findUser(username string) (string, string, string, string, string, string) {
	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	chk(err)

	defer db.Close()

	read, err := db.Query("SELECT * FROM users WHERE username='" + username + "';")
	chk(err)

	defer read.Close()

	if read.Next() {
		var a, b, c, d, e string

		err = read.Scan(&a, &b, &c, &d, &e)
		chk(err)

		return "Usuario encontrado: ", a, b, c, d, e
	} else {
		return "No se ha encontrado el usuario", "", "", "", "", ""
	}
}

//DELETE
func deleteUser(username string) string {
	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)
	chk(err)

	defer db.Close()

	delete, err := db.Query("DELETE FROM users WHERE username='" + username + "';")
	chk(err)

	defer delete.Close()

	return "Usuario borrado: " + username
}
