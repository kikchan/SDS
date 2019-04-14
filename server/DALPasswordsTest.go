package main

import (
	"fmt"
	"strconv"
)

//A group of miscelaneous tests
func DALPasswordsTest() {
	var code int
	var msg string

	fmt.Println("Create a password for an existing user")
	code, msg = createPassword("jose", "qwerty", "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Create a password for a non-existing user")
	code, msg = createPassword("jose", "qwerty", "juancarlos")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Search for an existing password of an existing user")
	code, msg = findPasswordByID("jose", 1)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Search for a non-existing password")
	code, msg = findPasswordByID("kiril", 1)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")
}
