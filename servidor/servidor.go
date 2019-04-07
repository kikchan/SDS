package main

import (
	"bufio"
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

func main() {
	var port = "8080"
	var code int
	var msg string

	/*
		CADUser function calls test
	*/
	code, msg = createUser("kiril", "123456", "Kiril", "Gaydarov", "kvg1@alu.ua.es")
	fmt.Println("(code: " + strconv.Itoa(code) + ", msg: " + msg + ")")

	code, msg = findUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", msg: " + msg + ")")

	code, msg = findUser("jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", msg: " + msg + ")")

	code, msg = updateUser("kiril", "654321", "kiril_gaydarov@gmail.com")
	fmt.Println("(code: " + strconv.Itoa(code) + ", msg: " + msg + ")")

	code, msg = findUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", msg: " + msg + ")")

	code, msg = deleteUser("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", msg: " + msg + ")")
	/*
		CADUser function calls test END
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
