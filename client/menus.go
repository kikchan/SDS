package main

import "fmt"

func menu(eleccion *int) {
	menu :=
		"\nWelcome to MasterPass\n" +
			"[ 1 ] Login\n" +
			"[ 2 ] Register\n" +
			"[ 3 ] Exit\n" +
			"Choose an option: "

	fmt.Print(menu)

	fmt.Scanln(eleccion)
}

func menuLogged(option *int, username string) {
	menuLogged :=
		"[ 1 ] Manage passwords\n" +
			"[ 2 ] Manage cards\n" +
			"[ 3 ] Manage notes\n" +
			"[ 4 ] Settings\n" +
			"[ 5 ] Logout\n" +
			"Choose an option: "

	fmt.Println(fmt.Sprintf("Welcome %s.\n", username))
	fmt.Print(menuLogged)
	fmt.Scanln(option)
}

func menuMngPasswords(option *int) {
	menuMngPasswords :=
		"[ 1 ] Add a password\n" +
			"[ 2 ] Show passwords\n" +
			"[ 3 ] Edit a password\n" +
			"[ 4 ] Delete a password\n" +
			"[ 5 ] Go back\n" +
			"Choose an option: "

	fmt.Print(menuMngPasswords)
	fmt.Scanln(option)
}

func menuMngCards(option *int) {
	menuMngCards :=
		"[ 1 ] Add a card\n" +
			"[ 2 ] Show cards\n" +
			"[ 3 ] Edit a card\n" +
			"[ 4 ] Delete a card\n" +
			"[ 5 ] Go back\n" +
			"Choose an option: "

	fmt.Print(menuMngCards)
	fmt.Scanln(option)
}

func menuMngNotes(option *int) {
	menuMngNotes :=
		"[ 1 ] Add a note\n" +
			"[ 2 ] Show notes\n" +
			"[ 3 ] Edit a note\n" +
			"[ 4 ] Delete a note\n" +
			"[ 5 ] Go back\n" +
			"Choose an option: "

	fmt.Print(menuMngNotes)
	fmt.Scanln(option)
}
