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

	. "github.com/logrusorgru/aurora"
)

func managePasswords(client *http.Client, username string) {
	clearScreen()

	//Creates a map of passwords
	passwords := make(map[int]passwordsData)

	//Request structure
	data := url.Values{}
	var option int

	//Set the "getUserPasswords" command
	data.Set("cmd", "getUserPasswords")

	//Set the username
	data.Set("username", username)

	//Send the request to the server
	r, err := client.PostForm(Server, data)
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

	switch option {
	case 1: //Add a password
		clearScreen()

		//Call the form to gather all the password data
		pd := addPassword()

		//Insert the new password into the map
		passwords[len(passwords)+1] = pd

		out, err := json.Marshal(passwords)
		if err != nil {
			panic(err)
		}

		data.Set("cmd", "modifyPasswords")
		data.Set("username", username)
		data.Set("passwords", encode64(out))

		r, err := client.PostForm(Server, data)
		chk(err)

		//Shows the response's body
		io.Copy(os.Stdout, r.Body)
		fmt.Println()

		return

	case 2: //List all passwords
		clearScreen()

		showPasswords(passwords, true)

		return

	case 3: //Edit a password
		clearScreen()

		if showPasswords(passwords, false) {
			fmt.Print("\n\nWhich password do you want to edit?: ")
			var index int
			fmt.Scanf("%d", &index)

			if index > 0 && index < len(passwords)+1 {
				//Call the form to gather all password's data
				pd := addPassword()

				//Replace the desired password with the newly generated one
				passwords[index] = pd

				out, err := json.Marshal(passwords)
				if err != nil {
					panic(err)
				}

				data.Set("cmd", "modifyPasswords")
				data.Set("username", username)
				data.Set("passwords", encode64(out))

				//Send the new data to the server so it can be stored
				r, err := client.PostForm(Server, data)
				chk(err)

				//Prints the server's response body
				io.Copy(os.Stdout, r.Body)
				fmt.Println()
			} else {
				fmt.Println(Red("The selected password doesn't exist"))

				time.Sleep(2 * time.Second)
			}
		}

		return

	case 4: //Delete password
		clearScreen()

		if showPasswords(passwords, false) {
			fmt.Print("\n\nWhich password do you want to delete?: ")
			var index int
			fmt.Scanf("%d", &index)

			if index > 0 && index < len(passwords)+1 {
				//Deletes the selected password from the map
				delete(passwords, index)

				out, err := json.Marshal(passwords)
				if err != nil {
					panic(err)
				}

				data.Set("cmd", "modifyPasswords")
				data.Set("username", username)
				data.Set("passwords", encode64(out))

				//Send the new data to the server so it can be stored
				r, err := client.PostForm(Server, data)
				chk(err)

				//Prints the server's response body
				io.Copy(os.Stdout, r.Body)
				fmt.Println()
			} else {
				fmt.Println(Red("The selected password doesn't exist"))

				time.Sleep(2 * time.Second)
			}
		}

		return

	case 5:
		logged(client, username)

	default:
		InvalidChoice()
		managePasswords(client, username)

	}
}

func manageCards(client *http.Client, username string) {
	clearScreen()

	//Creates a map of cards
	cards := make(map[int]cardsData)

	//Request structure
	data := url.Values{}
	var option int

	//Set the "getUserCards" command
	data.Set("cmd", "getUserCards")

	//Set the username
	data.Set("username", username)

	//Send the request to the server
	r, err := client.PostForm(Server, data)
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

		//Convert the response to a structure ot cards
		json.Unmarshal(decode64(m.Msg), &cards)
	}
	//------------------------------------------------------------------------

	menuMngCards(&option)

	switch option {
	case 1: //Add a card
		clearScreen()

		//Call the form to gather all the card data
		pd := addCard()

		//Insert the new card into the map
		cards[len(cards)+1] = pd

		out, err := json.Marshal(cards)
		if err != nil {
			panic(err)
		}

		data.Set("cmd", "modifyCards")
		data.Set("username", username)
		data.Set("cards", encode64(out))

		//Send the data to the server
		r, err := client.PostForm(Server, data)
		chk(err)

		//Shows the response's body
		io.Copy(os.Stdout, r.Body)
		fmt.Println()

		return

	case 2: //List all cards
		clearScreen()

		showCards(cards, true)

		return

	case 3: //Edit a card
		clearScreen()

		if showCards(cards, false) {
			fmt.Print("\n\nWhich card do you want to edit?: ")
			var index int
			fmt.Scanf("%d", &index)

			if index > 0 && index < len(cards)+1 {
				//Call the form to gather all card's data
				cd := addCard()

				//Replace the desired card with the newly generated one
				cards[index] = cd

				out, err := json.Marshal(cards)
				if err != nil {
					panic(err)
				}

				data.Set("cmd", "modifyCards")
				data.Set("username", username)
				data.Set("cards", encode64(out))

				//Send the new data to the server so it can be stored
				r, err := client.PostForm(Server, data)
				chk(err)

				//Prints the server's response body
				io.Copy(os.Stdout, r.Body)
				fmt.Println()
			} else {
				fmt.Println(Red("The selected card doesn't exist"))

				time.Sleep(2 * time.Second)
			}
		}

		return

	case 4: //Delete a card
		clearScreen()

		if showCards(cards, false) {
			fmt.Print("\n\nWhich card do you want to delete?: ")
			var index int
			fmt.Scanf("%d", &index)

			if index > 0 && index < len(cards)+1 {
				//Deletes the selected card from the map
				delete(cards, index)

				out, err := json.Marshal(cards)
				if err != nil {
					panic(err)
				}

				data.Set("cmd", "modifyCards")
				data.Set("username", username)
				data.Set("cards", encode64(out))

				//Send the new data to the server so it can be stored
				r, err := client.PostForm(Server, data)
				chk(err)

				//Prints the server's response body
				io.Copy(os.Stdout, r.Body)
				fmt.Println()
			} else {
				fmt.Println(Red("The selected card doesn't exist"))

				time.Sleep(2 * time.Second)
			}
		}

		return

	case 5:
		logged(client, username)

	default:
		InvalidChoice()
		manageCards(client, username)

	}
}

func manageNotes(client *http.Client, username string) {
	clearScreen()

	notas := make(map[int]notesData)

	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "getUserNotes") // comando (string)
	data.Set("username", username)

	r, err := client.PostForm(Server, data) // enviamos por POST
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

			r, err := client.PostForm(Server, data) // enviamos por POST
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

			r, err := client.PostForm(Server, data) // enviamos por POST
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

			r, err := client.PostForm(Server, data) // enviamos por POST
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
