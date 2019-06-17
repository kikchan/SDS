package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

//Default server address
var ServerIP = "https://localhost"
var ServerPort = "8080"
var Server = ServerIP + ":" + ServerPort

var keyData []byte

func main() {
	clearScreen()

	if len(os.Args) == 1 || len(os.Args) == 3 {
		var clear = false

		//Request structure
		var data = url.Values{}

		if len(os.Args) == 3 {
			ServerIP = "https://" + os.Args[1]
			ServerPort = os.Args[2]
			Server = ServerIP + ":" + ServerPort

			fmt.Println("Trying to establish connection with \"" + Server + "\" ...")
		} else if len(os.Args) == 1 {
			fmt.Println("Trying to establish local connection to default address: " + Server + " ...")
		}

		//Open an HTTP connection using TLS
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		//Create a special client that doesn't verify the self-signed certificate
		client := &http.Client{Transport: tr}

		data.Set("cmd", "open")
		data.Set("address", GetAddress())

		r, err := client.PostForm(Server, data)

		if err == nil && r.StatusCode == 200 {
			var response resp

			//Read the body from the response
			body, _ := ioutil.ReadAll(r.Body)

			//Creates a new JSON decoder
			dec := json.NewDecoder(strings.NewReader(string(body)))

			//Convert the body to a json structure
			dec.Decode(&response)
			fmt.Println(Green(response.Msg))

			//Wait 2 seconds so the user can read the message
			time.Sleep(1 * time.Second)

			for {
				//Clear the screen only the first time the program is executed
				if clear {
					clearScreen()
				} else {
					clear = true
				}

				//Show the public main menu
				publicMenu(client)
			}
		} else {
			fmt.Println(Red("Could not establish connection with requested server"))
		}
	} else {
		fmt.Println("Bad arguments. The correct sintax is: [programName] [server] [port]")
		fmt.Println("An example: \"go run *.go 127.0.0.1 8080\"")
	}
}

func publicMenu(client *http.Client) {
	clearScreen()

	//Request structure
	var data = url.Values{}

	//Response structure
	var m resp

	//User's choice
	var option int

	menu(&option)

	switch option {
	case 1:
		clearScreen()

		var username string
		var password string

		var tries int
		for tries = 5; tries >= 0; tries-- {
			login(&username, &password)

			//Hash the password with SHA512
			keyClient := sha512.Sum512([]byte(password))

			//Half of the password is used for the login (256 bits)
			keyLogin := keyClient[:32]

			//Store the second half of the password
			keyData = keyClient[32:64]

			//Set the "login" command
			data.Set("cmd", "login")

			//Set the username
			data.Set("user", username)

			//Set the user password (encoded)
			data.Set("pass", encode64(keyLogin))

			//Send the request to the server and get a response
			r, err := client.PostForm(Server, data)
			chk(err)

			//If connection is succsessfull, show the logged menu, otherwise, print a message
			if r.StatusCode == 200 {
				logged(client, username)
			} else {
				if tries == 1 {
					fmt.Println("Wrong password, " + strconv.Itoa(tries) + " try left.")
				} else {
					fmt.Println("Wrong password, " + strconv.Itoa(tries) + " tries left.")
				}
				fmt.Println()
			}

			//If there are no more tries left, make the user wait for 5 minutes
			if tries == 0 {
				fmt.Println("Please try again in 5 minutes")
				time.Sleep(5 * time.Minute)

				main()
			}
		}

	case 2:
		clearScreen()

		var username string
		var password string
		var name string
		var surname string
		var email string

		register(&username, &password, &name, &surname, &email)

		//Generate a pair of both private and public keys
		pkClient, err := rsa.GenerateKey(rand.Reader, 1024)
		chk(err)

		//Speeds up future operations
		pkClient.Precompute()

		//Parse the keys to a JSON structure
		pkJSON, err := json.Marshal(&pkClient)
		chk(err)

		//Extract the public key
		keyPub := pkClient.Public()

		//Parse it to JSON
		pubJSON, err := json.Marshal(&keyPub) // y codificamos con JSON
		chk(err)

		//Hash the password with SHA512
		keyClient := sha512.Sum512([]byte(password))

		//Half of the password is used for the login (256 bits)
		keyLogin := keyClient[:32]

		//The other half of the password is used for the data (256 bits)
		keyData := keyClient[32:64]

		//Create a user structure
		a := &userData{name, surname, email, encode64(encrypt(compress(pkJSON), keyData))}

		//Parse the user structure to JSON
		out, err := json.Marshal(a)
		if err != nil {
			panic(err)
		}

		//Set the "register" command
		data.Set("cmd", "register")

		//Set the username
		data.Set("user", username)

		//Set the user password
		data.Set("pass", encode64(keyLogin))

		//Set the user data
		data.Set("userData", encode64(out))

		//Set the compressed and encoded public key
		data.Set("pubkey", encode64(compress(pubJSON)))

		//Send the request to the server
		r, err := client.PostForm(Server, data)
		chk(err)

		//Read the body from the response
		body, _ := ioutil.ReadAll(r.Body)

		processResponse(body, &m)

		if m.Code == 1 {
			logged(client, username)
		}

	case 3:
		//Set the "close" command
		data.Set("cmd", "close")

		//Set the current client's IP
		data.Set("address", GetAddress())

		//Notify the server of the connection close
		client.PostForm(Server, data)

		//Exit the program
		fmt.Println("Goodbye!")
		os.Exit(0)

	default:
		InvalidChoice()
		publicMenu(client)
	}
}

func logged(client *http.Client, username string) {
	var option int

	clearScreen()

	menuLogged(&option, username)

	for {
		switch option {
		case 1:
			clearScreen()

			//Call the password manager function
			managePasswords(client, username)

		case 2:
			clearScreen()

			//Call the card manager function
			manageCards(client, username)

		case 3:
			clearScreen()

			//Call the note manager function
			manageNotes(client, username)

		case 4:
			clearScreen()

			//Call the user settings function
			userSettings(client, username)

		case 5:
			clearScreen()

			//Go back to previous menu
			publicMenu(client)

		default:
			InvalidChoice()
			logged(client, username)
		}
	}
}
