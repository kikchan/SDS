package main

import "fmt"

func menu(option *int) {
	menu :=
		"\nWelcome to MasterPass\n" +
			"[ 1 ] Login\n" +
			"[ 2 ] Register\n" +
			"[ 3 ] Exit\n" +
			"Choice: "

	fmt.Print(menu)
	fmt.Scanln(option)
}

func menuLogged(option *int, username string) {
	menuLogged :=
		"[ 1 ] Manage passwords\n" +
			"[ 2 ] Manage cards\n" +
			"[ 3 ] Manage notes\n" +
			"[ 4 ] User settings\n" +
			"[ 5 ] Logout\n" +
			"Choice: "

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
			"Choice: "

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
			"Choice: "

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
			"Choice: "

	fmt.Print(menuMngNotes)
	fmt.Scanln(option)
}

func menuUserSettings(option *int) {
	clearScreen()

	menuUserSettings :=
		"[ 1 ] View user data\n" +
			"[ 2 ] Change name\n" +
			"[ 3 ] Change surname\n" +
			"[ 4 ] Change email\n" +
			"[ 5 ] Delete account\n" +
			"[ 6 ] Go back\n" +
			"Choice: "

	fmt.Print(menuUserSettings)
	fmt.Scanln(option)
}
