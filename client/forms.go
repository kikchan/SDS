package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func login(username *string, password *string) {
	fmt.Println(Bold(Green("\nLogin")))
	fmt.Println("----------------------")

	fmt.Printf("Username: ")
	fmt.Scanln(username)

	fmt.Printf("Password: ")
	fmt.Scanln(password)
}

func register(username *string, password *string, name *string, surname *string, email *string) {
	fmt.Println(Bold(Green("\nRegister")))
	fmt.Println("----------------------")

	fmt.Print("Enter username: ")
	fmt.Scanln(username)

	fmt.Print("Enter password: ")
	fmt.Scanln(password)

	fmt.Print("Enter name: ")
	fmt.Scanln(name)

	fmt.Print("Enter surname: ")
	fmt.Scanln(surname)

	fmt.Print("Enter email: ")
	fmt.Scanln(email)
}
