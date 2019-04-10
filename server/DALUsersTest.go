package main

import (
	"fmt"
	"strconv"
)

//A group of miscelaneous tests
func DALUsersTest() {
	var code int
	var msg string

	fmt.Println("Create a user")
	code, msg = createUser("kiril", "123456", "Kiril", "Gaydarov", "kvg1@alu.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Duplicate same user, throws an error")
	code, msg = createUser("kiril", "123456", "Kiril", "Gaydarov", "kvg1@alu.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find the user that was just created")
	code, msg = findUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find a user that doesn't exist")
	code, msg = findUser("juan")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find another user that exists")
	code, msg = findUser("jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Update an existing user")
	code, msg = updateUser("kiril", "654321", "kiril_gaydarov@gmail.com")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Update a non existing user")
	code, msg = updateUser("kiril123", "654321", "kiril_gaydarov@gmail.com")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Delete an existing user")
	code, msg = deleteUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Delete a non-existing user")
	code, msg = deleteUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")
}
