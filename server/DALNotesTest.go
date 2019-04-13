package main

import (
	"fmt"
	"strconv"
)

//A group of miscelaneous tests
func DALNotesTest() {
	var code int
	var msg string

	fmt.Println("Create a note")
	code, msg = createNote("Hola, esto es una prueba!", "kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find a note")
	code, msg = findNoteByID("kiril", 11)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")
}
