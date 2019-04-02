package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//Función para crear un usuario
func createUser(username string, password string, name string, surname string) string {
	db, err := sql.Open("mysql", DB_Username+":"+DB_Password+"@"+DB_Protocol+"("+DB_IP+":"+DB_Port+")/"+DB_Name)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//insert, err := db.Query("INSERT INTO sds.users (username, password, name, surname) VALUES (sds, sds, SDS, SDS);")
	insert, err := db.Query("INSERT INTO users VALUES ('" + username + "', '" + password + "', '" + name + "', '" + surname + "');")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	return "Usuario creado: {" + username + ", " + name + " " + surname + "}"
}

//Función para borrar un usuario
func deleteUser(username string) string {
	db, err := sql.Open("mysql", "sds:sds@tcp(185.207.145.237:3306)/sds")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//insert, err := db.Query("INSERT INTO sds.users (username, password, name, surname) VALUES (sds, sds, SDS, SDS);")
	insert, err := db.Query("DELETE FROM users WHERE username='" + username + "';")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	return "Usuario borrado: {" + username + "}"
}
