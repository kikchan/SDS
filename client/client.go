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
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
)

//Default server address
var SERVER_IP string = "https://localhost"
var SERVER_PORT string = "8080"
var SERVER string = SERVER_IP + ":" + SERVER_PORT

func main() {
	if len(os.Args) == 1 || len(os.Args) == 3 {
		if len(os.Args) == 3 {
			SERVER_IP = os.Args[1]
			SERVER_PORT = os.Args[2]

			fmt.Println("Trying to establish connection with \"" + SERVER + "\"")
		} else if len(os.Args) == 1 {
			fmt.Println("Trying to establish local connection to default port: " + SERVER_PORT)
		}

		var eleccion int     //Declarar variable y tipo antes de escanear, esto es obligatorio
		data := url.Values{} // estructura para contener los valores

		/* creamos un cliente especial que no comprueba la validez de los certificados
		esto es necesario por que usamos certificados autofirmados (para pruebas) */
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		data.Set("cmd", "check")

		r, err := client.PostForm(SERVER, data)

		if err == nil {
			for {
				io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
				fmt.Println()

				menu(&eleccion)

				switch eleccion {
				case 1:
					clearScreen()

					var username string
					var password string
					data := url.Values{} // estructura para contener los valores

					var i int
					fmt.Println("Log in:")
					fmt.Println("-------------------")

					for i = 2; i >= 0; i-- {
						login(&username, &password)

						// hash con SHA512 de la contraseña
						keyClient := sha512.Sum512([]byte(password))
						keyLogin := keyClient[:32] // una mitad para el login (256 bits)

						data.Set("cmd", "login")             // comando (string)
						data.Set("user", username)           // usuario (string)
						data.Set("pass", encode64(keyLogin)) // contraseña (a base64 porque es []byte)
						r, err := client.PostForm(SERVER, data)
						chk(err)

						if r.StatusCode == 200 {
							logueado(client, username)
						} else {
							fmt.Println("CONTRASEÑA INVÁLIDA, te quedan", i, "intentos.")
							fmt.Println()
						}
						if i == 0 {
							fmt.Println("No podrás volver a intentarlo hasta dentro de 5 minutos.")
							time.Sleep(5 * time.Minute)
							fmt.Println("Ya puedes volver a intentarlo.")
							i = 3
						}
					}

				case 2:
					clearScreen()

					var username string
					var password string
					var name string
					var surname string
					var email string

					fmt.Println("Register:")
					fmt.Println("-------------------")
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

					r, err := client.PostForm(SERVER, data) // enviamos por POST
					chk(err)
					io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
					fmt.Println()

				case 3:
					clearScreen()

					fmt.Println("Goodbye!")
					return
				default:
					fmt.Println("Choose an option or press [Ctrl] + [C] to exit")
				}
			}
		} else {
			fmt.Println("Could not establish connection with requested server")
		}
	} else {
		fmt.Println("Bad arguments. The correct sintax is: [programName] [server] [port]")
		fmt.Println("An example: \"go run *.go https://localhost 8080\"")
	}
}

func login(username *string, password *string) {
	fmt.Printf("Username: ")
	fmt.Scanln(username)

	fmt.Printf("Password: ")
	fmt.Scanln(password)
}

func register(username *string, password *string, name *string, surname *string, email *string) {

	fmt.Println("Insert username:")
	fmt.Scanln(username)

	fmt.Println("Insert password:")
	fmt.Scanln(password)

	fmt.Println("Insert name:")
	fmt.Scanln(name)

	fmt.Println("Insert surname:")
	fmt.Scanln(surname)

	fmt.Println("Insert email:")
	fmt.Scanln(email)
}

func logueado(client *http.Client, username string) {
	var option int
	data := url.Values{} // estructura para contener los valores

	clearScreen()

	menuLogged(&option, username)

	for {
		switch option {
		case 1:
			managePasswords(client, username)

		case 2: //GestionCard
			manageCards(client, username)

		case 3: //GestionNote
			manageNotes(client, username)

		case 4:
			data.Set("cmd", "deleteUser")  // comando (string)
			data.Set("username", username) // usuario (string)

			r, err := client.PostForm(SERVER, data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()
			return

		case 5: //Share passwords

			return

		case 6:
			clearScreen()
			fmt.Println("Logged out")
			main()

		default:
			fmt.Println("Please choose a valid option")
		}
	}
}

func managePasswords(client *http.Client, username string) {
	contraseñas := make(map[int]passwordsData)
	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "Passwords") // comando (string)
	data.Set("username", username)

	r, err := client.PostForm(SERVER + "/cards", data) // enviamos por POST
	chk(err)

	//--------- Con esto recojo del servidor las tarjetas y las convierto al struct
	body, err := ioutil.ReadAll(r.Body)

	dec := json.NewDecoder(strings.NewReader(string(body)))

	for {
		var m resp
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(decode64(m.Msg), &contraseñas) // Con esto paso al map notas lo que recojo en el servidor
	}
	//------------------------------------------------------------------------

	var option int
	menuMngPasswords(&option)

	for {
		switch option {
		case 1: //addPassword
			newPassword := passwordsData{}

			data := url.Values{} // estructura para contener los valores

			data.Set("cmd", "modifyPasswords") // comando (string)

			fmt.Print("Inserte URL: ")
			var url string
			fmt.Scanf("%s", &url)

			fmt.Print("Inserte usuario: ")
			var user string
			fmt.Scanf("%s", &user)

			fmt.Print("¿Quieres generar una contraseña aleatoria?(s/n) ")
			var opcion string
			fmt.Scanf("%s", &opcion)

			var contraseña string
			if opcion == "s" {
				fmt.Print("Inserte longitud de la contraseña: ")
				var long int
				fmt.Scanf("%d", &long)

				fmt.Print("Inserte número de digitos de la contraseña: ")
				var numDigitos int
				fmt.Scanf("%d", &numDigitos)

				fmt.Print("Inserte número de simbolos de la contraseña: ")
				var numSimbolos int
				fmt.Scanf("%d", &numSimbolos)

				fmt.Print("¿Permitir mayúsculas y minusculas?(t/f): ")
				var upperLower bool
				fmt.Scanf("%t", &upperLower)

				fmt.Print("¿Repetir carácteres?(t/f): ")
				var repeatCharacers bool
				fmt.Scanf("%t", &repeatCharacers)
				// Generate a password that is 64 characters long with 10 digits, 10 symbols,
				// allowing upper and lower case letters, disallowing repeat characters.
				// upperLower = false es que permite
				contrasenyaa, err := password.Generate(long, numDigitos, numSimbolos, !upperLower, repeatCharacers)
				if err != nil {
					log.Fatal(err)
				}
				contraseña = contrasenyaa
			} else {
				fmt.Print("Introduce contraseña: ")
				fmt.Scanf("%s", &contraseña)
			}

			fmt.Printf("La contraseña generada es: ")
			fmt.Println(contraseña)

			newPassword.Username = user
			newPassword.Password = contraseña
			newPassword.Site = url
			newPassword.Modified = time.Now().String()

			contraseñas[len(contraseñas)+1] = newPassword

			out, err := json.Marshal(contraseñas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("passwords", encode64(out))

			r, err := client.PostForm(SERVER + "/passwords", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 2: //List password
			for k, v := range contraseñas {
				fmt.Println(k, "-. URL: ", v.Site, ", de: ", v.Username)
			}

			return

		case 3: //Modify password
			modifyPasswords := passwordsData{}
			data := url.Values{} // estructura para contener los valores

			for k, v := range contraseñas {
				fmt.Println(k, "-. URL: ", v.Site, ", de: ", v.Username)
			}

			data.Set("cmd", "modifyPasswords") // comando (string)

			fmt.Print("¿Que contraseña quieres editar?(num) ")
			var index int
			fmt.Scanf("%d", &index)

			fmt.Print("Inserte URL: ")
			var url string
			fmt.Scanf("%s", &url)

			fmt.Print("Inserte usuario: ")
			var user string
			fmt.Scanf("%s", &user)

			fmt.Print("¿Quieres generar una contraseña aleatoria?(s/n) ")
			var opcion string
			fmt.Scanf("%s", &opcion)

			var contraseña string
			if opcion == "s" {
				fmt.Print("Inserte longitud de la contraseña: ")
				var long int
				fmt.Scanf("%d", &long)

				fmt.Print("Inserte número de digitos de la contraseña: ")
				var numDigitos int
				fmt.Scanf("%d", &numDigitos)

				fmt.Print("Inserte número de simbolos de la contraseña: ")
				var numSimbolos int
				fmt.Scanf("%d", &numSimbolos)

				fmt.Print("¿Permitir mayúsculas y minusculas?(t/f): ")
				var upperLower bool
				fmt.Scanf("%t", &upperLower)

				fmt.Print("¿Repetir carácteres?(t/f): ")
				var repeatCharacers bool
				fmt.Scanf("%t", &repeatCharacers)
				// Generate a password that is 64 characters long with 10 digits, 10 symbols,
				// allowing upper and lower case letters, disallowing repeat characters.
				// upperLower = false es que permite
				contrasenyaa, err := password.Generate(long, numDigitos, numSimbolos, !upperLower, repeatCharacers)
				if err != nil {
					log.Fatal(err)
				}
				contraseña = contrasenyaa
			} else {
				fmt.Print("Introduce contraseña: ")
				fmt.Scanf("%s", &contraseña)
			}

			fmt.Printf("La contraseña generada es: ")
			fmt.Println(contraseña)

			modifyPasswords.Username = user
			modifyPasswords.Password = contraseña
			modifyPasswords.Site = url
			modifyPasswords.Modified = time.Now().String()

			contraseñas[index] = modifyPasswords
			out, err := json.Marshal(contraseñas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("passwords", encode64(out))

			r, err := client.PostForm(SERVER + "/passwords", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()
			return

		case 4: //Delete password
			data := url.Values{} // estructura para contener los valores

			for k, v := range contraseñas {
				fmt.Println(k, "-. URL: ", v.Site, ", de: ", v.Username)
			}

			data.Set("cmd", "modifyPasswords") // comando (string)

			fmt.Print("¿Que contraseña quieres borrar?(num) ")
			var index int
			fmt.Scanf("%d", &index)

			delete(contraseñas, index)

			// KIRIL, AQUI HABRIA QUE HACER ALGO PARA QUE SE VUELVAN A COLOCAR LOS INDEX DEL MAP (YA QUE AL BORRAR SE QUEDA UNO SUELTO)

			out, err := json.Marshal(contraseñas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("passwords", encode64(out))

			r, err := client.PostForm(SERVER + "/passwords", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 5:
			return

		default:
			fmt.Println("No has soptionado una opcion correcta.")

		}
	}
}

func manageCards(client *http.Client, username string) {
	tarjetas := make(map[int]cardsData)

	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "Cards") // comando (string)
	data.Set("username", username)

	r, err := client.PostForm(SERVER + "/cards", data) // enviamos por POST
	chk(err)

	//--------- Con esto recojo del servidor las tarjetas y las convierto al struct
	body, err := ioutil.ReadAll(r.Body)

	dec := json.NewDecoder(strings.NewReader(string(body)))

	for {
		var m resp
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(decode64(m.Msg), &tarjetas) // Con esto paso al map notas lo que recojo en el servidor
	}
	//------------------------------------------------------------------------

	var option int
	menuMngCards(&option)

	for {
		switch option {
		case 1: //addCard
			newCard := cardsData{}

			data := url.Values{} // estructura para contener los valores

			data.Set("cmd", "modifyCards") // comando (string)

			fmt.Print("Inserte propietario de la tarjeta: ")
			var owner string
			fmt.Scanf("%s", &owner)

			fmt.Print("Inserte número de la tarjeta: ")
			var pan string
			fmt.Scanf("%s", &pan)

			fmt.Print("Inserte CCV: ")
			var ccv string
			fmt.Scanf("%s", &ccv)

			fmt.Print("Inserte mes de caducidad: ")
			var month string
			fmt.Scanf("%s", &month)

			fmt.Print("Inserte año de caducidad: ")
			var year string
			fmt.Scanf("%s", &year)

			newCard.Pan = pan
			newCard.Owner = owner
			newCard.Ccv = ccv
			newCard.Expiry = month + "-" + year

			tarjetas[len(tarjetas)+1] = newCard

			out, err := json.Marshal(tarjetas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("cards", encode64(out))

			r, err := client.PostForm(SERVER + "/cards", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 2: //List cards
			for k, v := range tarjetas {
				fmt.Println(k, "-. Número: ", v.Pan, ", de: ", v.Owner)
			}

			return

		case 3: //Modify card
			modifyCard := cardsData{}
			data := url.Values{} // estructura para contener los valores

			for k, v := range tarjetas {
				fmt.Println(k, "-. Número: ", v.Pan, ", de: ", v.Owner)
			}

			data.Set("cmd", "modifyCards") // comando (string)

			fmt.Print("¿Que tarjeta quieres editar?(num) ")
			var index int
			fmt.Scanf("%d", &index)

			fmt.Print("Inserte propietario de la tarjeta: ")
			var owner string
			fmt.Scanf("%s", &owner)

			fmt.Print("Inserte número de la tarjeta: ")
			var pan string
			fmt.Scanf("%s", &pan)

			fmt.Print("Inserte CCV: ")
			var ccv string
			fmt.Scanf("%s", &ccv)

			fmt.Print("Inserte mes de caducidad: ")
			var month string
			fmt.Scanf("%s", &month)

			fmt.Print("Inserte año de caducidad: ")
			var year string
			fmt.Scanf("%s", &year)

			modifyCard.Pan = pan
			modifyCard.Owner = owner
			modifyCard.Ccv = ccv
			modifyCard.Expiry = month + "-" + year

			tarjetas[index] = modifyCard
			out, err := json.Marshal(tarjetas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("cards", encode64(out))

			r, err := client.PostForm(SERVER + "/cards", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()
			return

		case 4: //Delete card
			data := url.Values{} // estructura para contener los valores

			for k, v := range tarjetas {
				fmt.Println(k, "-. Número: ", v.Pan, ", de: ", v.Owner)
			}

			data.Set("cmd", "modifyCards") // comando (string)

			fmt.Print("¿Que tarjeta quieres borrar?(num) ")
			var index int
			fmt.Scanf("%d", &index)

			delete(tarjetas, index)

			// KIRIL, AQUI HABRIA QUE HACER ALGO PARA QUE SE VUELVAN A COLOCAR LOS INDEX DEL MAP (YA QUE AL BORRAR SE QUEDA UNO SUELTO)

			out, err := json.Marshal(tarjetas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("cards", encode64(out))

			r, err := client.PostForm(SERVER + "/cards", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 5:
			return

		default:
			fmt.Println("No has soptionado una opcion correcta.")

		}
	}
}

func manageNotes(client *http.Client, username string) {
	notas := make(map[int]notesData)

	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "Notes") // comando (string)
	data.Set("username", username)

	r, err := client.PostForm(SERVER + "/notes", data) // enviamos por POST
	chk(err)

	//--------- Con esto recojo del servidor las notas y las convierto al struct
	body, err := ioutil.ReadAll(r.Body)

	dec := json.NewDecoder(strings.NewReader(string(body)))

	for {
		var m resp
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(decode64(m.Msg), &notas) // Con esto paso al map notas lo que recojo en el servidor
	}
	//------------------------------------------------------------------------

	var option int
	menuMngNotes(&option)

	for {
		switch option {
		case 1: //add note
			newNota := notesData{}

			data := url.Values{} // estructura para contener los valores

			data.Set("cmd", "modifyNotes") // comando (string)

			fmt.Print("Inserte nota: ")
			var text string
			fmt.Scanf("%s", &text)

			fmt.Print("Inserte fecha: ")
			var date string
			fmt.Scanf("%s", &date)

			newNota.Text = text
			newNota.Date = date

			notas[len(notas)+1] = newNota
			out, err := json.Marshal(notas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("notas", encode64(out))

			r, err := client.PostForm(SERVER + "/notes", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 2: //List notes
			for k, v := range notas {
				fmt.Println(k, "-. Texto: ", v.Text, ", fecha: ", v.Date)
			}

			return

		case 3: //Modify note
			modifyNota := notesData{}
			data := url.Values{} // estructura para contener los valores

			for k, v := range notas {
				fmt.Println(k, "-. Texto: ", v.Text, ", fecha: ", v.Date)
			}

			data.Set("cmd", "modifyNotes") // comando (string)

			fmt.Print("¿Que nota quieres editar?(num) ")
			var index int
			fmt.Scanf("%d", &index)

			fmt.Print("Inserte nueva nota: ")
			var text string
			fmt.Scanf("%s", &text)

			fmt.Print("Inserte nueva fecha: ")
			var date string
			fmt.Scanf("%s", &date)

			modifyNota.Text = text
			modifyNota.Date = date

			notas[index] = modifyNota
			out, err := json.Marshal(notas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("notas", encode64(out))

			r, err := client.PostForm(SERVER + "/notes", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()
			return

		case 4: //Delete note
			data := url.Values{} // estructura para contener los valores

			for k, v := range notas {
				fmt.Println(k, "-. Texto: ", v.Text, ", fecha: ", v.Date)
			}

			data.Set("cmd", "modifyNotes") // comando (string)

			fmt.Print("¿Que nota quieres borrar?(num) ")
			var index int
			fmt.Scanf("%d", &index)

			delete(notas, index)

			// KIRIL, AQUI HABRIA QUE HACER ALGO PARA QUE SE VUELVAN A COLOCAR LOS INDEX DEL MAP (YA QUE AL BORRAR SE QUEDA UNO SUELTO)

			out, err := json.Marshal(notas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("notas", encode64(out))

			r, err := client.PostForm(SERVER + "/notes", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 5:
			return

		default:
			fmt.Println("No has seleccionado una opcion correcta.")

		}
	}
}
