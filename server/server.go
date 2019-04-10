package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

//Database connection variables
var DB_IP string = "185.207.145.237"
var DB_Port string = "3306"
var DB_Protocol string = "tcp"
var DB_Name string = "sds"
var DB_Username string = "sds"
var DB_Password string = "sds"

//Error check function
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func main() {
	var port = "8080"
	var code int
	var msg string

	/*
		DALUsers function calls test
	*/
	//Create a user
	code, msg = createUser("kiril", "123456", "Kiril", "Gaydarov", "kvg1@alu.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Duplicate same user, throws an error
	code, msg = createUser("kiril", "123456", "Kiril", "Gaydarov", "kvg1@alu.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Find the user that was just created
	code, msg = findUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Find a user that doesn't exist
	code, msg = findUser("juan")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Find another user that exists
	code, msg = findUser("jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Update an existing user
	code, msg = updateUser("kiril", "654321", "kiril_gaydarov@gmail.com")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Update a non existing user
	code, msg = updateUser("kiril123", "654321", "kiril_gaydarov@gmail.com")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	code, msg = findUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Delete an existing user
	code, msg = deleteUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Delete a non-existing user
	code, msg = deleteUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")
	/*
		DALUsers function calls test END
	*/

	/*
		DALCards function calls test
	*/
	//Create a card
	code, msg = createCard("123456879832158", "111", 05, 2030, "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Duplicate the previous card, should thrown an error
	code, msg = createCard("123456879832158", "111", 05, 2030, "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Find the card by PAN
	code, msg = findCardByPAN("jose", "123456879832158")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Find the card by ID
	code, msg = findCardByID("jose", 68)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Find a non-existing card
	code, msg = findCardByID("jose", 1)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Find all cards for user "jose"
	code, msg = getUserCards("jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Update an existing card
	code, msg = updateCard("123456879832158", "155", 04, 2019, "jose", "123456879832158")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Update a non-existing card
	code, msg = updateCard("56469547632115", "155", 04, 2019, "jose", "15")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Delete an existing card
	code, msg = deleteCard("123456879832158", "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	//Update a non-existing card
	code, msg = deleteCard("0123", "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")
	/*
		DALCards function calls test END
	*/

	if len(os.Args) == 2 {
		port = os.Args[1]
		fmt.Println("Server awaiting connections from port: " + port)
	} else {
		fmt.Println("Server awaiting connections from port: " + port + " (default)")
	}

	//Server is in listening mode
	ln, err := net.Listen("tcp", "localhost:"+port)
	chk(err)

	defer ln.Close()

	//Infinite loop
	for {
		//Accept every single user request
		conn, err := ln.Accept()
		chk(err)

		//Launch a concurrent lambda function
		go func() {
			//Gets the user's port
			_, port, err := net.SplitHostPort(conn.RemoteAddr().String())
			chk(err)

			fmt.Println("Connection: ", conn.LocalAddr(), " <--> ", conn.RemoteAddr())

			scanner := bufio.NewScanner(conn)

			//Scans the connection and reads the message
			for scanner.Scan() {
				//Print the user's message
				fmt.Println("Client[", port, "]: ", scanner.Text())

				//Send "ACK" to client
				fmt.Fprintln(conn, "ack: ", scanner.Text())
			}

			conn.Close()
			fmt.Println("Closed[", port, "]")
		}()
	}
}
