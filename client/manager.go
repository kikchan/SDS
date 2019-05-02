package main

import (
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

func managePasswords(client *http.Client, username string) {
	clearScreen()

	passwords := make(map[int]passwordsData)

	//Request structure
	data := url.Values{}
	var option int

	//Set the "getUserPasswords" command
	data.Set("cmd", "getUserPasswords")

	//Set the username
	data.Set("username", username)

	//Send the request to the server
	r, err := client.PostForm(Server+"/cards", data)
	chk(err)

	//Retrieve the response's body
	body, err := ioutil.ReadAll(r.Body)

	//Create a new JSON decoder
	dec := json.NewDecoder(strings.NewReader(string(body)))

	for {
		var m resp

		//Decode the server's response
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		//Convert the response to a structure of passwords
		json.Unmarshal(decode64(m.Msg), &passwords)
	}
	//------------------------------------------------------------------------

	menuMngPasswords(&option)

	for {
		switch option {
		case 1:
			data.Set("cmd", "modifyPasswords")

			newPassword := addPassword()

			passwords[len(passwords)+1] = newPassword

			out, err := json.Marshal(passwords)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("passwords", encode64(out))

			r, err := client.PostForm(Server+"/passwords", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 2: //List password
			for k, v := range passwords {
				fmt.Println(k, "-. URL: ", v.Site, ", de: ", v.Username)
			}

			return

		case 3: //Modify password
			modifyPasswords := passwordsData{}

			for k, v := range passwords {
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
			var random string
			fmt.Scanf("%s", &random)

			var contraseña string
			if random == "s" {
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

			passwords[index] = modifyPasswords
			out, err := json.Marshal(passwords)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("passwords", encode64(out))

			r, err := client.PostForm(Server+"/passwords", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()
			return

		case 4: //Delete password
			for k, v := range passwords {
				fmt.Println(k, "-. URL: ", v.Site, ", de: ", v.Username)
			}

			data.Set("cmd", "modifyPasswords") // comando (string)

			fmt.Print("¿Que contraseña quieres borrar?(num) ")
			var index int
			fmt.Scanf("%d", &index)

			delete(passwords, index)

			// KIRIL, AQUI HABRIA QUE HACER ALGO PARA QUE SE VUELVAN A COLOCAR LOS INDEX DEL MAP (YA QUE AL BORRAR SE QUEDA UNO SUELTO)

			out, err := json.Marshal(passwords)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("passwords", encode64(out))

			r, err := client.PostForm(Server+"/passwords", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 5:
			logged(client, username)

		default:
			InvalidChoice()
			managePasswords(client, username)

		}
	}
}

func manageCards(client *http.Client, username string) {
	clearScreen()

	tarjetas := make(map[int]cardsData)

	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "getUserCards") // comando (string)
	data.Set("username", username)

	r, err := client.PostForm(Server+"/cards", data) // enviamos por POST
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

			r, err := client.PostForm(Server+"/cards", data) // enviamos por POST
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

			r, err := client.PostForm(Server+"/cards", data) // enviamos por POST
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

			r, err := client.PostForm(Server+"/cards", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 5:
			logged(client, username)

		default:
			InvalidChoice()
			manageCards(client, username)

		}
	}
}

func manageNotes(client *http.Client, username string) {
	clearScreen()

	notas := make(map[int]notesData)

	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "getUserNotes") // comando (string)
	data.Set("username", username)

	r, err := client.PostForm(Server+"/notes", data) // enviamos por POST
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

			r, err := client.PostForm(Server+"/notes", data) // enviamos por POST
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

			r, err := client.PostForm(Server+"/notes", data) // enviamos por POST
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

			r, err := client.PostForm(Server+"/notes", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 5:
			logged(client, username)

		default:
			InvalidChoice()
			manageNotes(client, username)

		}
	}
}

func userSettings(client *http.Client, username string) {
	clearScreen()

	var option int

	data := url.Values{} //Request structure

	menuUserSettings(&option)

	switch option {
	case 1: //View user data

	case 2: //Change name

	case 3: //Change surname

	case 4: //Change email

	case 5: //Delete account
		data.Set("cmd", "deleteUser")  // comando (string)
		data.Set("username", username) // usuario (string)

		r, err := client.PostForm(Server, data) // enviamos por POST
		chk(err)
		io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
		fmt.Println()
		return

	case 6: //Go back
		logged(client, username)

	default:
		InvalidChoice()
		userSettings(client, username)
	}
}
