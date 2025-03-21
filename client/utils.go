package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

//Error checking function
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

//Invalid option message
func InvalidChoice() {
	fmt.Println(Red("Please choose a valid option"))
	time.Sleep(2 * time.Second)
}

//A screen cleaner. Only works on Unix based systems
func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

//Encrypts the given data using AES
func encrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)+16)    //Reserve space at the begginig of the array
	rand.Read(out[:16])                 //Reads the output array
	blk, err := aes.NewCipher(key)      //Creates a new AES cipher using the key
	chk(err)                            //Check if there's any error
	ctr := cipher.NewCTR(blk, out[:16]) //Cipher using CTR
	ctr.XORKeyStream(out[16:], data)    //Encrypt the data

	return
}

//Generates a random AES key with a size of 32 bytes
func generateAESkey() []byte {
	key := make([]byte, 32)
	rand.Read(key)

	return key
}

func decrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)-16)
	blk, err := aes.NewCipher(key)
	chk(err)
	ctr := cipher.NewCTR(blk, data[:16])
	ctr.XORKeyStream(out, data[16:])
	return
}

//Compress the given data
func compress(data []byte) []byte {
	var b bytes.Buffer      //Define a variable to store the compressed data
	w := zlib.NewWriter(&b) //Creates a writer to compress on b
	w.Write(data)           //Write data
	w.Close()               //Close the writter
	return b.Bytes()
}

//Decompress the given data
func decompress(data []byte) []byte {
	var b bytes.Buffer                              //Contains the uncompressed data
	r, err := zlib.NewReader(bytes.NewReader(data)) //A reader that uncompresses
	chk(err)                                        //Check if there's any error
	io.Copy(&b, r)                                  //Copy form the reader to the variable
	r.Close()                                       //Closes the reader
	return b.Bytes()
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

//Get the IP of the client machine
func GetAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	return conn.LocalAddr().String()
}

//Prints a colored menu
func coloredMenu(title string) {
	fmt.Println(Green("--------------------------"))
	fmt.Println(Green("|"), Red("\t"+title+"\t"), Green("|"))
	fmt.Println(Green("--------------------------\n"))
}

//Receives the response's body then parses it to JSON and checks the returned Code from the server
func processResponse(body []byte, m *resp) int {
	//Creates a new JSON decoder
	dec := json.NewDecoder(strings.NewReader(string(body)))

	//Convert the body to a json structure
	dec.Decode(&m)

	fmt.Println()

	if m.Code == 1 {
		fmt.Println(Green(m.Msg))
	} else {
		fmt.Println(Red(m.Msg))
	}

	fmt.Print("Press any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	return m.Code
}

//Receives the response's body then parses it to JSON then returns it as an array
func convertResponseToArrayOfUsers(body []byte, m *resp) []user {
	//Creates a new JSON decoder
	dec := json.NewDecoder(strings.NewReader(string(body)))

	//Convert the body to a json structure
	dec.Decode(&m)

	//Array of users
	var users []user

	//Array of users splitted
	arrayOfUsers := strings.Split(m.Msg, "###")

	for i := range arrayOfUsers {
		//Each user's data gets splitted in another array of 2 columns
		arrayOfUserData := strings.Split(arrayOfUsers[i], "##")

		var user user
		user.Username = arrayOfUserData[0]
		user.PubKey = arrayOfUserData[1]

		users = append(users, user)
	}

	return users
}

//Receives the response's body then parses it to JSON then returns it as an array
func decryptResponseToArrayOfPasswords(body []byte, m *resp, pKey string) []passwordsData {
	//Creates a new JSON decoder
	dec := json.NewDecoder(strings.NewReader(string(body)))

	//Convert the body to a json structure
	dec.Decode(&m)

	//Array of passwords
	var passwords []passwordsData

	if m.Code == 1 {
		//Array of passwords splitted
		arrayOfPasswords := strings.Split(m.Msg, "###")

		for i := range arrayOfPasswords {
			//Each password's data gets splitted in another array of 2 columns
			passStruct := strings.Split(arrayOfPasswords[i], "##")

			//Get the user's private key using the second half of his login password (keyData)
			var privateKey rsa.PrivateKey
			err := json.Unmarshal(decompress(decrypt(decode64(pKey), keyData)), &privateKey)
			chk(err)

			//Get the password's AES key that's going to be used to decrypt it
			AESkey, err := rsa.DecryptPKCS1v15(rand.Reader, &privateKey, decode64(passStruct[1]))
			chk(err)

			//Decrypt and parse to a struct the password
			var pd passwordsData
			err = json.Unmarshal(decrypt(decode64(passStruct[0]), AESkey), &pd)

			//Insert the password into the array
			passwords = append(passwords, pd)
		}
	}

	return passwords
}

//An error message printing function
func invalidIndex(typeOfField string) {
	fmt.Println(Red("The selected " + typeOfField + " doesn't exist"))
	fmt.Println("Press any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

//Checks if the index exists in the given map of passwords
func passwordIndexExists(index int, array map[int]passwordsData) bool {
	var exists = false

	for i := range array {
		if index == i {
			exists = true
		}
	}

	return exists
}

//Checks if the index exists in the given map of cards
func cardIndexExists(index int, array map[int]cardsData) bool {
	var exists = false

	for i := range array {
		if index == i {
			exists = true
		}
	}

	return exists
}

//Checks if the index exists in the given map of notes
func noteIndexExists(index int, array map[int]notesData) bool {
	var exists = false

	for i := range array {
		if index == i {
			exists = true
		}
	}

	return exists
}
