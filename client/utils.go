package main

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net"
	"os"
	"os/exec"
)

//Error checking function
func chk(e error) {
	if e != nil {
		panic(e)
	}
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

//Compress the given data
func compress(data []byte) []byte {
	var b bytes.Buffer      //Define a variable to store the compressed data
	w := zlib.NewWriter(&b) //Creates a writer to compress on b
	w.Write(data)           //Write data
	w.Close()               //Close the writter
	return b.Bytes()        // devolvemos los datos comprimidos
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
