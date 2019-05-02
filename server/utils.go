package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/logrusorgru/aurora"
)

//Error checking function
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

//Convert from []byte to string
func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) //Encode the string
}

//Convert from string to []byte
func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) //Decode the string
	chk(err)                                     //Check if there's any error
	return b
}

//Send response to client. The response contains a code and a message
func response(w http.ResponseWriter, code int, msg string) {
	if code != 1 {
		w.WriteHeader(400)
	}

	r := resp{Code: code, Msg: msg} //Build the response
	rJSON, err := json.Marshal(&r)  //Convert it to JSON
	chk(err)                        //Check for errors
	w.Write(rJSON)                  //Send the response
}

//Prints out to the console the current timestamp
func consoleTimeStamp() {
	fmt.Print(Bold(Green("[" + time.Now().Format("2006-01-02 15:04:05") + "]: ")))
}
