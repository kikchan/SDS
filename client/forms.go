package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	"github.com/sethvargo/go-password/password"
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

func addPassword() passwordsData {
	var pd passwordsData
	var random string

	pd.Modified = time.Now().Format("2006-01-02 15:04:05")

	fmt.Print("Enter the URL of the password's site: ")
	fmt.Scanln(&pd.Site)

	fmt.Print("Enter your username: ")
	fmt.Scanln(&pd.Username)

	fmt.Print("Would you like to generate a random password for it? (y/n): ")
	fmt.Scanln(&random)

	if random == "y" {
		var size, nDigits, nSymbols int
		var choice string
		var upperLower = false
		var repeat = false

		fmt.Print("Size of the password: ")
		fmt.Scanln(&size)

		fmt.Print("Number of digits: ")
		fmt.Scanln(&nDigits)

		fmt.Print("Number of symbols: ")
		fmt.Scanln(&nSymbols)

		fmt.Print("Allow upper and lowercase letters? (y/n): ")
		fmt.Scanln(&choice)

		if choice == "y" {
			upperLower = true
		}

		fmt.Print("Repeat characters? (y/n): ")
		fmt.Scanln(&choice)

		if choice == "y" {
			repeat = true
		}

		// Generate a password that is 64 characters long with 10 digits, 10 symbols,
		// allowing upper and lower case letters, disallowing repeat characters.
		pass, err := password.Generate(size, nDigits, nSymbols, !upperLower, repeat)
		if err != nil {
			log.Fatal(err)
		}

		pd.Password = pass

		fmt.Print(Bold(Red("Showing the generated password for")), Underline(Bold(White("5 seconds!"))), "\n\n")
		fmt.Println(pd.Password)

		time.Sleep(5 * time.Second)
	} else {
		fmt.Print("Enter your password: ")
		fmt.Scanln(&pd.Password)
	}

	return pd
}

func addCard() cardsData {
	var cd cardsData
	in := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the name of the owner: ")
	line, _ := in.ReadString('\n')
	cd.Owner = strings.TrimSuffix(line, "\n")

	fmt.Print("Enter the card's identification number (PAN): ")
	fmt.Scanln(&cd.Pan)

	fmt.Print("Enter the card's secred number (CCV): ")
	fmt.Scanln(&cd.Ccv)

	fmt.Print("Enter the card's expiry date (E.g. 01/20): ")
	fmt.Scanln(&cd.Expiry)

	return cd
}

func addNote() notesData {
	var nd notesData
	in := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the note's text: ")
	line, _ := in.ReadString('\n')
	nd.Text = line

	nd.Date = time.Now().Format("2006-01-02 15:04:05")

	return nd
}

func deleteUser() bool {
	fmt.Print("This will", Red("erase"), "all your stores passwords, cards and notes. Are you sure? (y/n): ")
	var choice string

	fmt.Scanln(&choice)

	if choice == "y" {
		return true
	} else {
		return false
	}
}
