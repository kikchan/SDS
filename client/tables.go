package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	. "github.com/logrusorgru/aurora"
)

func showPasswords(passwords map[int]passwordsData, displayPass bool) bool {
	var foundAnyPasswords bool

	if len(passwords) > 0 {
		if displayPass {
			fmt.Print(Bold(Red("Showing passwords for")), Underline(Bold(White("5 seconds!"))), "\n\n")
			fmt.Println(" #\t | Site\t\t| Username\t| Password\t| Last modified")
			fmt.Println("-----------------------------------------------------------------------------")
		} else {
			fmt.Println(" #\t | Site\t\t| Username\t| Last modified")
			fmt.Println("-------------------------------------------------------------")
		}

		for k, v := range passwords {
			fmt.Print("[" + strconv.Itoa(k) + "]\t | ")

			if len(v.Site) >= 5 {
				fmt.Print(Blue(v.Site), "\t| ")
			} else {
				fmt.Print(Blue(v.Site), "\t\t| ")
			}

			if len(v.Username) >= 6 {
				fmt.Print(Bold(Green(v.Username)), "\t| ")
			} else {
				fmt.Print(Bold(Green(v.Username)), "\t\t| ")
			}

			if displayPass {
				if len(v.Password) >= 6 {
					fmt.Print(Bold(Red(v.Password)), "\t| ")
				} else {
					fmt.Print(Bold(Red(v.Password)), "\t\t| ")
				}
			}

			fmt.Println(v.Modified)
		}

		foundAnyPasswords = true

		if displayPass {
			time.Sleep(5 * time.Second)
		}
	} else {
		fmt.Println(Red("No passwords to show"))

		foundAnyPasswords = false

		fmt.Println("Press any key to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	return foundAnyPasswords
}

func showCards(cards map[int]cardsData, displayCCV bool) bool {
	var foundAnyCards bool

	if len(cards) > 0 {
		if displayCCV {
			fmt.Print(Bold(Red("Showing cards for")), Underline(Bold(White("5 seconds!"))), "\n\n")
			fmt.Println(" #\t | Owner\t\t\t| PAN\t\t\t| CCV\t| Expiry date")
			fmt.Println("-------------------------------------------------------------------------------------")
		} else {
			fmt.Println(" #\t | Owner\t\t\t| PAN\t\t\t| Expiry date")
			fmt.Println("-----------------------------------------------------------------------------")
		}

		for k, v := range cards {
			fmt.Print("[" + strconv.Itoa(k) + "]\t | ")

			if len(v.Owner) >= 20 {
				fmt.Print(Blue(v.Owner), "\t| ")
			} else {
				fmt.Print(Blue(v.Owner), "\t\t| ")
			}

			fmt.Print(Bold(Green(v.Pan)), "\t| ")

			if displayCCV {
				fmt.Print(Bold(Red(v.Ccv)), "\t| ")
			}

			fmt.Println(v.Expiry)
		}

		foundAnyCards = true

		if displayCCV {
			time.Sleep(5 * time.Second)
		}
	} else {
		fmt.Println(Red("No cards to show"))

		foundAnyCards = false

		fmt.Println("Press any key to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	return foundAnyCards
}

func showNotes(notes map[int]notesData, displayMsg bool) bool {
	var foundAnyNotes bool

	if len(notes) > 0 {
		fmt.Println(" #\t | Date \t\t | Text")
		fmt.Println("---------------------------------------------------------------------------")

		for k, v := range notes {
			fmt.Print("[" + strconv.Itoa(k) + "]\t | ")
			fmt.Print(v.Date, "\t | ")
			fmt.Println(Blue(v.Text))
		}

		foundAnyNotes = true

		if displayMsg {
			fmt.Print("\nPress any key to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	} else {
		fmt.Println(Red("No notes to show"))

		foundAnyNotes = false

		fmt.Print("Press any key to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	return foundAnyNotes
}

func showUserData(userData userData, username string) {
	var choice string

	fmt.Println(Blue("Name:\t"), userData.Name)
	fmt.Println(Blue("Surname:"), userData.Surname)
	fmt.Println(Blue("Email:\t"), userData.Email)

	fmt.Print("\nDo you want to see the private key? (y/n): ")
	fmt.Scanln(&choice)

	if choice == "y" {
		fmt.Println(Blue("Private key:\n"), userData.PrivateKey)

		fmt.Print("\nPress any key to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
