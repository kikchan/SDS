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
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	//"net/http"
	//"crypto/tls"
	"io"

	"github.com/sethvargo/go-password/password"
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

// respuesta del servidor
type resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
}

// función para cifrar (con AES en este caso), adjunta el IV al principio
func encrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)+16)    // reservamos espacio para el IV al principio
	rand.Read(out[:16])                 // generamos el IV
	blk, err := aes.NewCipher(key)      // cifrador en bloque (AES), usa key
	chk(err)                            // comprobamos el error
	ctr := cipher.NewCTR(blk, out[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out[16:], data)    // ciframos los datos
	return
}

// función para comprimir
func compress(data []byte) []byte {
	var b bytes.Buffer      // b contendrá los datos comprimidos (tamaño variable)
	w := zlib.NewWriter(&b) // escritor que comprime sobre b
	w.Write(data)           // escribimos los datos
	w.Close()               // cerramos el escritor (buffering)
	return b.Bytes()        // devolvemos los datos comprimidos
}

// función para codificar de []bytes a string (Base64)
func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

func menu(eleccion *int) {
	menu :=
		`
		Bienvenido
		[ 1 ] Login
		[ 2 ] Register
		¿Qué prefieres?

	`

	fmt.Print(menu)

	fmt.Scanln(eleccion)
}

func login(username *string, password *string) {
	fmt.Printf("Insert username: ")
	fmt.Scanln(username)

	fmt.Printf("Insert password: ")
	fmt.Scanln(password)
}

func register(username *string, password *string, name *string, surname *string, email *string) {

	fmt.Println("Insert username:")
	fmt.Scanln(username)

	fmt.Println("Insert password:")
	fmt.Scanln(password)

	fmt.Println("Insert name:")
	fmt.Scanln(name)

	fmt.Println("Insert surname:")
	fmt.Scanln(surname)

	fmt.Println("Insert email:")
	fmt.Scanln(email)
}

func main() {
	var puerto = "8080"

	if len(os.Args) == 2 {
		puerto = os.Args[1]
		fmt.Println("Intentando conectar al puerto: " + puerto)
	} else {
		fmt.Println("Intentando conectar al puerto: " + puerto + " (por defecto)")
	}

	var eleccion int //Declarar variable y tipo antes de escanear, esto es obligatorio

	/* creamos un cliente especial que no comprueba la validez de los certificados
	esto es necesario por que usamos certificados autofirmados (para pruebas) */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// generamos un par de claves (privada, pública) para el servidor
	pkClient, err := rsa.GenerateKey(rand.Reader, 1024)
	chk(err)
	pkClient.Precompute() // aceleramos su uso con un precálculo

	pkJSON, err := json.Marshal(&pkClient) // codificamos con JSON
	chk(err)

	keyPub := pkClient.Public()           // extraemos la clave pública por separado
	pubJSON, err := json.Marshal(&keyPub) // y codificamos con JSON
	chk(err)

	for {
		menu(&eleccion)

		switch eleccion {
		case 1:
			var username string
			var password string

			fmt.Println("Iniciar sesión:")
			login(&username, &password)

			// hash con SHA512 de la contraseña
			keyClient := sha512.Sum512([]byte(password))
			keyLogin := keyClient[:32] // una mitad para el login (256 bits)

			// ** ejemplo de login
			data := url.Values{}
			data.Set("cmd", "login")             // comando (string)
			data.Set("user", username)           // usuario (string)
			data.Set("pass", encode64(keyLogin)) // contraseña (a base64 porque es []byte)
			r, err := client.PostForm("https://localhost:8080", data)
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)

			if r.StatusCode == 200 {
				logueado(&data)
			}
		case 2:
			var username string
			var password string
			var name string
			var surname string
			var email string

			fmt.Println("Registrar usuario:")
			register(&username, &password, &name, &surname, &email)

			// hash con SHA512 de la contraseña
			keyClient := sha512.Sum512([]byte(password))
			keyLogin := keyClient[:32]  // una mitad para el login (256 bits)
			keyData := keyClient[32:64] // la otra para los datos (256 bits)

			// ** ejemplo de registro
			data := url.Values{}                 // estructura para contener los valores
			data.Set("cmd", "register")          // comando (string)
			data.Set("user", username)           // usuario (string)
			data.Set("pass", encode64(keyLogin)) // "contraseña" a base64
			data.Set("name", name)
			data.Set("surname", surname)
			data.Set("email", email)

			// comprimimos y codificamos la clave pública
			data.Set("pubkey", encode64(compress(pubJSON)))

			// comprimimos, ciframos y codificamos la clave privada
			data.Set("prikey", encode64(encrypt(compress(pkJSON), keyData)))

			r, err := client.PostForm("https://localhost:8080", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()
		default:
			fmt.Println("No prefieres ninguno de ellos")
		}
	}
}

func menuLogueado(eleccion *int) {
	menuLogueado :=
		`
		Bienvenido
		[ 1 ] Añadir contraseña de sitio web
		[ 2 ] Ver contraseña de sitio web
		[ 3 ] Añadir tarjeta de crédito
		[ 4 ] Ver tarjetas de crédito
		[ 5 ] Eliminar tu usuario
		[ 6 ] Cerrar sesión
		¿Qué prefieres?

	`
	fmt.Print(menuLogueado)
	fmt.Scanln(eleccion)
}

func logueado(data *url.Values) {
	var eleccion int
	menuLogueado(&eleccion)

	for {

		switch eleccion {
		case 1: //Add password
			fmt.Print("Inserte URL: ")
			var url string
			fmt.Scanf("%s", &url)

			fmt.Print("Inserte longitud de la contraseña: ")
			var long int
			fmt.Scanf("%d", &long)

			fmt.Print("Inserte número de digitos de la contraseña: ")
			var numDigitos int
			fmt.Scanf("%d", &numDigitos)

			fmt.Print("Inserte número de simbolos de la contraseña: ")
			var numSimbolos int
			fmt.Scanf("%d", &numSimbolos)

			fmt.Print("¿Permitir mayúsculas y minusculas?(t/f): ")
			var upperLower bool
			fmt.Scanf("%t", &upperLower)

			fmt.Print("¿Repetir carácteres?(t/f): ")
			var repeatCharacers bool
			fmt.Scanf("%t", &repeatCharacers)
			// Generate a password that is 64 characters long with 10 digits, 10 symbols,
			// allowing upper and lower case letters, disallowing repeat characters.
			// upperLower = false es que permite
			contraseña, err := password.Generate(long, numDigitos, numSimbolos, !upperLower, repeatCharacers)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("La contraseña generada es: ")
			fmt.Println(contraseña)

		case 2: //View password

		case 3:

		case 4:

		case 5:

		case 6:
			fmt.Println("Hasta la vista.")
			return
		default:
			fmt.Println("No has seleccionado una opcion correcta.")

		}
	}
}
