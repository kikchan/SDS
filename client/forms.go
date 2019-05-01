package main

import "fmt"

func login(username *string, password *string) {
	fmt.Printf("Username: ")
	fmt.Scanln(username)

	fmt.Printf("Password: ")
	fmt.Scanln(password)
}

func register(username *string, password *string, name *string, surname *string, email *string) {

	fmt.Print("Insert username: ")
	fmt.Scanln(username)

	fmt.Print("Insert password: ")
	fmt.Scanln(password)

	fmt.Print("Insert name: ")
	fmt.Scanln(name)

	fmt.Print("Insert surname: ")
	fmt.Scanln(surname)

	fmt.Print("Insert email: ")
	fmt.Scanln(email)
}
