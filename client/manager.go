package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func managePasswords(client *http.Client, username string) {
	clearScreen()

	//Creates a map of passwords
	passwords := make(map[int]passwordsData)

	//Request structure
	data := url.Values{}

	//Response structure
	var m resp

	//User's choice
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

		//Read the body from the response
		body, _ := ioutil.ReadAll(r.Body)

		processResponse(body, &m)

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

			if index > 0 {
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

				//Read the body from the response
				body, _ := ioutil.ReadAll(r.Body)

				processResponse(body, &m)
			} else {
				invalidIndex("password")
			}
		}

		return

	case 4: //Delete a password
		clearScreen()

		if showPasswords(passwords, false) {
			fmt.Print("\n\nWhich password do you want to delete?: ")
			var index int
			fmt.Scanf("%d", &index)

			if index > 0 {
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

				//Read the body from the response
				body, _ := ioutil.ReadAll(r.Body)

				processResponse(body, &m)
			} else {
				invalidIndex("password")
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

	//Response structure
	var m resp

	//User's choice
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

		//Read the body from the response
		body, _ := ioutil.ReadAll(r.Body)

		processResponse(body, &m)

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

			if index > 0 {
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

				//Read the body from the response
				body, _ := ioutil.ReadAll(r.Body)

				processResponse(body, &m)
			} else {
				invalidIndex("card")
			}
		}

		return

	case 4: //Delete a card
		clearScreen()

		if showCards(cards, false) {
			fmt.Print("\n\nWhich card do you want to delete?: ")
			var index int
			fmt.Scanf("%d", &index)

			if index > 0 {
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

				//Read the body from the response
				body, _ := ioutil.ReadAll(r.Body)

				processResponse(body, &m)
			} else {
				invalidIndex("card")
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

	//Creates a map of notes
	notes := make(map[int]notesData)

	//Request structure
	data := url.Values{}

	//Response structure
	var m resp

	//User's choice
	var option int

	//Set the "getUserNotes" command
	data.Set("cmd", "getUserNotes")

	//Set the username
	data.Set("username", username)

	//Set the request to the server
	r, err := client.PostForm(Server, data)
	chk(err)

	//Retrieve the response's body
	body, err := ioutil.ReadAll(r.Body)

	//Create a new JSON decoder
	dec := json.NewDecoder(strings.NewReader(string(body)))

	for {
		//Decode the server's response
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		//Convert the response to a structure of notes
		json.Unmarshal(decode64(m.Msg), &notes)
	}
	//------------------------------------------------------------------------

	menuMngNotes(&option)

	switch option {
	case 1: //Add a note
		clearScreen()

		//Call the form to gather all the note data
		nd := addNote()

		//Insert the new note into the map
		notes[len(notes)+1] = nd

		out, err := json.Marshal(notes)
		if err != nil {
			panic(err)
		}

		data.Set("cmd", "modifyNotes")
		data.Set("username", username)
		data.Set("notes", encode64(out))

		//Send the request to the server
		r, err := client.PostForm(Server, data)
		chk(err)

		//Read the body from the response
		body, _ := ioutil.ReadAll(r.Body)

		processResponse(body, &m)

		return

	case 2: //List all notes
		clearScreen()

		showNotes(notes, true)

		return

	case 3: //Edit a note
		clearScreen()

		if showNotes(notes, false) {
			fmt.Print("\nWhich note do you want to edit?: ")
			var index int
			fmt.Scanf("%d", &index)

			if index > 0 {
				//Call the form to gather all note's data
				nd := addNote()

				//Replace the desired note with the newly generated one
				notes[index] = nd

				out, err := json.Marshal(notes)
				if err != nil {
					panic(err)
				}

				data.Set("cmd", "modifyNotes")
				data.Set("username", username)
				data.Set("notes", encode64(out))

				//Send the new data to the server so it can be stored
				r, err := client.PostForm(Server, data)
				chk(err)

				//Read the body from the response
				body, _ := ioutil.ReadAll(r.Body)

				processResponse(body, &m)
			} else {
				invalidIndex("note")
			}
		}

		return

	case 4: //Delete a note
		clearScreen()

		if showNotes(notes, false) {
			fmt.Print("\nWhich note do you want to delete?: ")
			var index int
			fmt.Scanf("%d", &index)

			if index > 0 {
				//Deletes the selected note from the map
				delete(notes, index)

				out, err := json.Marshal(notes)
				if err != nil {
					panic(err)
				}

				data.Set("cmd", "modifyNotes")
				data.Set("username", username)
				data.Set("notes", encode64(out))

				//Send the new data to the server so it can be stored
				r, err := client.PostForm(Server, data)
				chk(err)

				//Read the body from the response
				body, _ := ioutil.ReadAll(r.Body)

				processResponse(body, &m)
			} else {
				invalidIndex("note")
			}
		}

		return

	case 5:
		logged(client, username)

	default:
		InvalidChoice()
		manageNotes(client, username)

	}
}

func userSettings(client *http.Client, username string) {
	clearScreen()

	//Request structure
	data := url.Values{}

	//Response structure
	var m resp

	//Set the "readUser" command
	data.Set("cmd", "readUser")

	//Set the username
	data.Set("username", username)

	//Send the request to the server
	r, err := client.PostForm(Server, data)
	chk(err)

	//Retrieve the response's body
	body, err := ioutil.ReadAll(r.Body)

	//Create a new JSON decoder
	dec := json.NewDecoder(strings.NewReader(string(body)))

	var usuario userData
	for {
		//Decode the server's response
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		//Convert the response to a structure of passwords
		json.Unmarshal(decode64(m.Msg), &usuario)
	}
	//------------------------------------------------------------------------

	//User's choice
	var option int

	menuUserSettings(&option)

	switch option {
	case 1: //View user data
		clearScreen()

		showUserData(usuario, username)

		return

	case 2: //Change name

		var newName string
		fmt.Printf("Enter new name: ")
		fmt.Scanf("%s", &newName)

		usuario.Name = newName

		//Request structure
		data := url.Values{}

		//Response structure
		var m resp

		//Set the "updateUser" command
		data.Set("cmd", "updateUser")

		out, err := json.Marshal(usuario)
		if err != nil {
			panic(err)
		}

		//Set the username
		data.Set("username", username)
		data.Set("data", encode64(out))

		//Send the request to the server
		r, err := client.PostForm(Server, data)
		chk(err)

		//Read the body from the response
		body, _ := ioutil.ReadAll(r.Body)

		processResponse(body, &m)

		return

	case 3: //Change surname

		var newSurname string
		fmt.Printf("Enter new surname: ")
		fmt.Scanf("%s", &newSurname)

		usuario.Surname = newSurname

		//Request structure
		data := url.Values{}

		//Response structure
		var m resp

		//Set the "updateUser" command
		data.Set("cmd", "updateUser")

		out, err := json.Marshal(usuario)
		if err != nil {
			panic(err)
		}

		//Set the username
		data.Set("username", username)
		data.Set("data", encode64(out))

		//Send the request to the server
		r, err := client.PostForm(Server, data)
		chk(err)

		//Read the body from the response
		body, _ := ioutil.ReadAll(r.Body)

		processResponse(body, &m)

		return

	case 4: //Change email

		var newEmail string
		fmt.Printf("Enter new email: ")
		fmt.Scanf("%s", &newEmail)

		usuario.Email = newEmail

		//Request structure
		data := url.Values{}

		//Response structure
		var m resp

		//Set the "updateUser" command
		data.Set("cmd", "updateUser")

		out, err := json.Marshal(usuario)
		if err != nil {
			panic(err)
		}

		//Set the username
		data.Set("username", username)
		data.Set("data", encode64(out))

		//Send the request to the server
		r, err := client.PostForm(Server, data)
		chk(err)

		//Read the body from the response
		body, _ := ioutil.ReadAll(r.Body)

		processResponse(body, &m)

		return

	case 5: //Delete account
		clearScreen()

		if deleteUser() {
			data.Set("cmd", "deleteUser")
			data.Set("username", username)

			//Send the new data to the server so the user gets deleted
			r, err := client.PostForm(Server, data)
			chk(err)

			body, _ := ioutil.ReadAll(r.Body)

			processResponse(body, &m)

			publicMenu(client)
		}

		return

	case 6: //Go back
		logged(client, username)

	default:
		InvalidChoice()
		userSettings(client, username)
	}
}
