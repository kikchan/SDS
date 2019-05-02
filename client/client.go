package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
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

func main() {
	clearScreen()

	if len(os.Args) == 1 || len(os.Args) == 3 {
		var clear = false
		var data = url.Values{} //Request structure

		if len(os.Args) == 3 {
			ServerIP = os.Args[1]
			ServerPort = os.Args[2]

			fmt.Println("Trying to establish connection with \"" + Server + "\" ...")
		} else if len(os.Args) == 1 {
			fmt.Println("Trying to establish local connection to default port: " + ServerPort + " ...")
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		//Create a special client that doesn't verify the self-signed certificate
		client := &http.Client{Transport: tr}

		data.Set("cmd", "check")
		data.Set("address", GetAddress())

		r, err := client.PostForm(Server, data)

		if err == nil && r.StatusCode == 200 {
			var response resp
			body, _ := ioutil.ReadAll(r.Body)
			dec := json.NewDecoder(strings.NewReader(string(body)))

			dec.Decode(&response)
			fmt.Println(Green(response.Msg))

			time.Sleep(1 * time.Second) //Wait 2 seconds so the user can read the message

			for {
				if clear {
					clearScreen()
				} else {
					clear = true
				}

				publicMenu(client)
			}
		} else {
			fmt.Println(Red("Could not establish connection with requested server"))
		}
	} else {
		fmt.Println("Bad arguments. The correct sintax is: [programName] [server] [port]")
		fmt.Println("An example: \"go run *.go https://localhost 8080\"")
	}
}

func publicMenu(client *http.Client) {
	clearScreen()

	var option int
	data := url.Values{} //Request structure

	menu(&option)

	switch option {
	case 1:
		var username string
		var password string

		var tries int
		for tries = 5; tries >= 0; tries-- {
			login(&username, &password)

			// hash con SHA512 de la contraseña
			keyClient := sha512.Sum512([]byte(password))
			keyLogin := keyClient[:32] // una mitad para el login (256 bits)

			data.Set("cmd", "login")             // comando (string)
			data.Set("user", username)           // usuario (string)
			data.Set("pass", encode64(keyLogin)) // contraseña (a base64 porque es []byte)

			r, err := client.PostForm(Server, data)
			chk(err)

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

			if tries == 0 {
				fmt.Println("Please try again in 5 minutes")
				time.Sleep(5 * time.Minute)

				main()
			}
		}

	case 2:
		var username string
		var password string
		var name string
		var surname string
		var email string

		register(&username, &password, &name, &surname, &email)

		// generamos un par de claves (privada, pública) para el servidor
		pkClient, err := rsa.GenerateKey(rand.Reader, 1024)
		chk(err)
		pkClient.Precompute() // aceleramos su uso con un precálculo

		pkJSON, err := json.Marshal(&pkClient) // codificamos con JSON
		chk(err)

		keyPub := pkClient.Public()           // extraemos la clave pública por separado
		pubJSON, err := json.Marshal(&keyPub) // y codificamos con JSON
		chk(err)

		// hash con SHA512 de la contraseña
		keyClient := sha512.Sum512([]byte(password))
		keyLogin := keyClient[:32]  // una mitad para el login (256 bits)
		keyData := keyClient[32:64] // la otra para los datos (256 bits)

		a := &userData{name, surname, email, encode64(encrypt(compress(pkJSON), keyData))}

		out, err := json.Marshal(a)
		if err != nil {
			panic(err)
		}

		data.Set("cmd", "register")          // comando (string)
		data.Set("user", username)           // usuario (string)
		data.Set("pass", encode64(keyLogin)) // "contraseña" a base64
		data.Set("userData", encode64(out))

		// comprimimos y codificamos la clave pública
		data.Set("pubkey", encode64(compress(pubJSON)))

		// comprimimos, ciframos y codificamos la clave privada
		//data.Set("prikey", encode64(encrypt(compress(pkJSON), keyData)))

		r, err := client.PostForm(Server, data) // enviamos por POST
		chk(err)
		io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
		fmt.Println()

	case 3:
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
			managePasswords(client, username)

		case 2:
			clearScreen()
			manageCards(client, username)

		case 3:
			clearScreen()
			manageNotes(client, username)

		case 4:
			clearScreen()
			userSettings(client, username)

		case 5:
			clearScreen()
			publicMenu(client)

			return

		default:
			InvalidChoice()
			logged(client, username)
		}
	}
}
