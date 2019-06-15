package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/logrusorgru/aurora"
	"golang.org/x/crypto/scrypt"
)

//Database connection variables
var DB_IP = "185.207.145.237"
var DB_Port = "3306"
var DB_Protocol = "tcp"
var DB_Name = "sds2"
var DB_Username = "sdsAppClient"
var DB_Password = "qwerty123456"

//Default parameters
var Port = "8080"

func main() {
	if len(os.Args) == 2 {
		Port = os.Args[1]
		fmt.Println("Awaiting connections through port:", Cyan(Port))
	} else {
		fmt.Println("Awaiting connections through", Magenta("default"), "port:", Yellow(Port))
	}

	//Asign a global handler to "/" route
	http.HandleFunc("/", handler)

	//Set the server to await connections from the selected port
	chk(http.ListenAndServeTLS(":"+Port, "cert.pem", "key.pem", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	//The request must be always parsed
	req.ParseForm()

	//Set the header to "text/plain"
	w.Header().Set("Content-Type", "text/plain")

	//Switch between the commands received from the client
	switch req.Form.Get("cmd") {
	case "login":
		//Find the user
		code, msg, user := findUser(req.Form.Get("user"))

		if code == 1 {
			//Find its passwords
			password := decode64(req.Form.Get("pass"))

			//Calculate the hash given its password
			hash, _ := scrypt.Key(password, user.Salt, 16384, 8, 1, 32)

			//Compare the calculated hash and the stored one
			if bytes.Compare(user.Hash, hash) == 0 {
				response(w, 1, "Logged in")
			} else {
				response(w, 0, "Wrong credentials")
			}
		} else {
			response(w, code, msg)
		}

		return

	case "register":
		//Find the user
		code, msg, _ := findUser(req.Form.Get("user"))

		if code == -3 {
			//Get user data
			u := user{}
			u.Username = req.Form.Get("user")
			u.Data = req.Form.Get("userData")
			u.PubKey = req.Form.Get("pubkey")

			//Make a 128 bits Salt
			u.Salt = make([]byte, 16)
			rand.Read(u.Salt)

			//Decode the user password
			password := decode64(req.Form.Get("pass"))

			//Calculate the hash given the user's password
			u.Hash, _ = scrypt.Key(password, u.Salt, 16384, 8, 1, 32)

			//Save the new user to the database with his public key, data, salt and hash
			code, msg = createUser(u.Username, u.PubKey, encode64(u.Hash), encode64(u.Salt), u.Data)

			response(w, code, msg)
		} else {
			response(w, 0, "The user already exists")
		}

	case "readUser":

		code, msg, user := findUser(req.Form.Get("username"))

		if code == 1 {
			response(w, 1, user.Data)
		} else {
			response(w, 0, msg)
		}

		return

	case "updateUser":

		code, msg := updateDataUser(req.Form.Get("username"), req.Form.Get("data"))

		if code == 1 {
			response(w, 1, msg)
		} else {
			response(w, 0, msg)
		}

		return

	case "deleteUser":
		//Find the user
		code, msg, username := findUser(req.Form.Get("username"))

		if code == 1 {
			//Try to delete the user
			code, msg := deleteUser(username.Username)

			response(w, code, msg)
		} else {
			response(w, code, msg)
		}

		return

	case "getUserPasswords":
		//Query the database to retrieve the user's passwords
		code, msg := getUserPasswords(req.Form.Get("username"))

		response(w, code, msg)

		return

	case "modifyPasswords":
		//Store in a variable the username and passwords to modify
		username := req.Form.Get("username")
		passwords := req.Form.Get("passwords")

		//Query the database to update the passwords
		code, msg := updatePassword(username, passwords)

		response(w, code, msg)

		return

	case "getUserCards":
		//Query the database to retrieve the user's cards
		code, msg := getUserCards(req.Form.Get("username"))

		response(w, code, msg)

		return

	case "modifyCards":
		//Store in a variable the username and cards to modify
		username := req.Form.Get("username")
		cards := req.Form.Get("cards")

		//Query the database to update the cards
		code, msg := updateCard(username, cards)

		response(w, code, msg)

		return

	case "getUserNotes":
		//Query the database to retreive the user's notes
		code, msg := getUserNotes(req.Form.Get("username"))

		response(w, code, msg)

		return

	case "modifyNotes":
		//Store in a variable the username and notes to modify
		username := req.Form.Get("username")
		notes := req.Form.Get("notes")

		//Query the database to update the notes
		code, msg := updateNote(username, notes)

		response(w, code, msg)

		return

	case "showAvailableUsers":
		//Store in a variable the username and notes to modify
		username := req.Form.Get("username")

		//Query the database to update the notes
		code, msg, users := getAllUsersExceptTheGivenOne(username)

		if len(users) > 0 {
			msg = ""

			for i, element := range users {
				msg += element.Username + "##" + element.PubKey

				if i != len(users)-1 {
					msg += "###"
				}
			}
		}

		response(w, code, msg)

		return

	case "open":
		response(w, 1, "Connection established!")

		consoleTimeStamp()
		fmt.Println(Yellow(req.Form.Get("address")), Green("opened"), BrightBlue("a new connection"))

	case "close":
		consoleTimeStamp()
		fmt.Println(Yellow(req.Form.Get("address")), Red("closed"), BrightBlue("the existing connection"))

	default:
		response(w, 0, "Please choose a valid option")
	}
}
