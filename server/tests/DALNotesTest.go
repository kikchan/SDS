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

	fmt.Println("Delete a note")
	code, msg = deleteNote(16, "kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Edit a note")
	code, msg = updateNote(11, "Esto parece funcionar bien...", "kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Get all user notes")
	code, msg = getUserNotes("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")
}
