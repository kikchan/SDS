package main

import (
	"fmt"
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

		time.Sleep(2 * time.Second)
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
				fmt.Print(Blue(v.Owner), "\t\t\t| ")
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

		time.Sleep(2 * time.Second)
	}

	return foundAnyCards
}
