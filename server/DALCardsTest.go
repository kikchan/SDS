package main

import (
	"fmt"
	"strconv"
)

//A group of miscelaneous tests
func DALCardsTest() {
	var code int
	var msg string

	fmt.Println("Create a card")
	code, msg = createCard("123456879832158", "111", 05, 2030, "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Duplicate the previous card, should thrown an error")
	code, msg = createCard("123456879832158", "111", 05, 2030, "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find the card by PAN")
	code, msg = findCardByPAN("jose", "123456879832158")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find the card by ID")
	code, msg = findCardByID("jose", 68)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find a non-existing card")
	code, msg = findCardByID("jose", 1)
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find all cards for user \"jose\"")
	code, msg = getUserCards("jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Find all cards for user \"kiril\"")
	code, msg = getUserCards("kiril")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Update an existing card")
	code, msg = updateCard("123456879832158", "155", 04, 2019, "jose", "123456879832158")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Update a non-existing card")
	code, msg = updateCard("56469547632115", "155", 04, 2019, "jose", "15")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Delete an existing card")
	code, msg = deleteCard("123456879832158", "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Delete a non-existing card")
	code, msg = deleteCard("0123", "jose")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")

	fmt.Println("Delete a non-existing card of a non-existing user")
	code, msg = deleteCard("4444", "emma")
	fmt.Println("(code: " + strconv.Itoa(code) + ", \tmsg: " + msg + ")")
}
