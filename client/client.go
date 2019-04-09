/*

Este programa demuestra una arquitectura cliente servidor sencilla.
El cliente envía líneas desde la entrada estandar y el servidor le devuelve un reconomiento de llegada (acknowledge).
El servidor es concurrente, siendo capaz de manejar múltiples clientes simultáneamente.
Las entradas se procesan mediante un scanner (bufio).

ejemplos de uso:

go run cnx.go srv

go run cnx.go cli

*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"crypto/md5"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	//"net/http"
	//"crypto/tls"
	"io"
)

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		//panic(e)
		fmt.Println("No se ha podido establecer una conexión con el servidor a través del puerto indicado")
		os.Exit(1)
	}
}

// ejemplo de tipo para un usuario
/*type user struct {
	username string            // nombre de usuario
	password string		   // Contraseña
	name string			   // Nombre
	surname string		   // Apelllidos
	email string		   // email (quitar)
	Hash []byte            // hash de la contraseña
	Salt []byte            // sal para la contraseña
	Data map[string]string // datos adicionales del usuario
}*/

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func menu(){
	menu := 
	`
		Bienvenido
		[ 1 ] Login
		[ 2 ] Register
		¿Qué prefieres?

	`

	fmt.Print(menu)

	var eleccion int //Declarar variable y tipo antes de escanear, esto es obligatorio
	fmt.Scanln(&eleccion)

	switch eleccion{
		case 1:
			fmt.Println("Iniciar sesión:")
			login()
		case 2:
			fmt.Println("Registrar usuario:")
			register()
		default:
			fmt.Println("No prefieres ninguno de ellos")
	}
}

func login(){
	var username string
	var password string

	fmt.Printf("Insert username: ")
	fmt.Scanln(&username)

	fmt.Printf("Insert password: ")
	fmt.Scanln(&password)

	ciphertext := encrypt([]byte(username), password)
	fmt.Printf("Encrypted: %x\n", ciphertext)
	plaintext := decrypt(ciphertext, password)
	fmt.Printf("Decrypted: %s\n", plaintext)

}

func register(){
	var username string
	var surname string
	var password string
	var creditCard string	
	
	fmt.Println("Insert username:")
	fmt.Scanln(&username)
	
	fmt.Println("Insert surname:")
	fmt.Scanln(&surname)

	fmt.Println("Insert password:")
	fmt.Scanln(&password)
	
	fmt.Println("Insert credit card:")
	fmt.Scanln(&creditCard)


	ciphertext := encrypt([]byte(username), password)
	fmt.Printf("Encrypted: %x\n", ciphertext)
	plaintext := decrypt(ciphertext, password)
	fmt.Printf("Decrypted: %s\n", plaintext)

}

func main() {
	var puerto = "8080"

	if len(os.Args) == 2 {
		puerto = os.Args[1]
		fmt.Println("Intentando conectar al puerto: " + puerto)
	} else {
		fmt.Println("Intentando conectar al puerto: " + puerto + " (por defecto)")
	}

	conn, err := net.Dial("tcp", "localhost:"+puerto) // llamamos al servidor
	chk(err)
	defer conn.Close() // es importante cerrar la conexión al finalizar

	fmt.Println("Entrando en modo cliente...")
	fmt.Println("conectado a ", conn.RemoteAddr())

	menu()

	keyscan := bufio.NewScanner(os.Stdin) // scanner para la entrada estándar (teclado)
	netscan := bufio.NewScanner(conn)     // scanner para la conexión (datos desde el servidor)

	for keyscan.Scan() { // escaneamos la entrada
		fmt.Fprintln(conn, keyscan.Text())         // enviamos la entrada al servidor
		netscan.Scan()                             // escaneamos la conexión
		fmt.Println("servidor: ", netscan.Text()) // mostramos mensaje desde el servidor
	}

	/*
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	r, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	chk(err)
	io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	fmt.Println()
	*/
}
