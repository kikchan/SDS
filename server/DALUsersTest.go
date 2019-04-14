package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"

	"golang.org/x/crypto/scrypt"
)

//A group of miscelaneous tests
func DALUsersTest() {
	var code int
	var msg string

	//MODIFY
	salt := make([]byte, 16) // sal (16 bytes == 128 bits)
	rand.Read(salt)          // la sal es aleatoria
	hash, _ := scrypt.Key([]byte("123456"), salt, 16384, 8, 1, 32)

	fmt.Println("Create a user")
	code, msg = createUser("kiril", "123456", hex.EncodeToString(hash), hex.EncodeToString(salt), "Kiril", "Gaydarov", "kvg1@alu.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Duplicate same user, throws an error")
	code, msg = createUser("kiril", "123456", hex.EncodeToString(hash), hex.EncodeToString(salt), "Kiril", "Gaydarov", "kvg1@alu.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find the user that was just created")
	code, msg, _ = findUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find a user that doesn't exist")
	code, msg, _ = findUser("juan")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find another user that exists")
	code, msg, _ = findUser("jose")
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
