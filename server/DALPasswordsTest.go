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
	code, msg = createPassword("jose", "qwerty", "jose", "https://www.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Create a password for a non-existing user")
	code, msg = createPassword("jose", "qwerty", "juancarlos", "https://www.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Search for an existing password of an existing user")
	code, msg = findPasswordByID("jose", 18)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Search for a non-existing password")
	code, msg = findPasswordByID("kiril", 1)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Get user passwords")
	code, msg = getUserPasswords("jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Get passwords for a non-existing user")
	code, msg = getUserPasswords("asdfasdfasdfasdf")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Get user passwords for site https://www.ua.es")
	code, msg = getPasswordsBySite("jose", "https://www.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Update password 19 to jose1 and 123456")
	code, msg = updatePassword(19, "jose1", "123456", "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Delete last password for jose")
	code, msg = deletePassword(28, "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")
}
