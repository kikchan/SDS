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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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

type notesData struct {
	Date string
	Text string
}

type userData struct {
	Name    string
	Surname string
	Email   string
}

// ejemplo de tipo para un usuario
type user struct {
	username string            // nombre de usuario
	password string            // Contraseña
	Hash     []byte            // hash de la contraseña
	Salt     []byte            // sal para la contraseña
	Data     map[string]string // datos adicionales del usuario
}

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

// función para decodificar de string a []bytes (Base64)
func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	chk(err)                                     // comprobamos el error
	return b                                     // devolvemos los datos originales
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
				logueado(client, username)
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

			a := &userData{name, surname, email}

			out, err := json.Marshal(a)
			if err != nil {
				panic(err)
			}

			// ** ejemplo de registro
			data := url.Values{}                 // estructura para contener los valores
			data.Set("cmd", "register")          // comando (string)
			data.Set("user", username)           // usuario (string)
			data.Set("pass", encode64(keyLogin)) // "contraseña" a base64
			data.Set("userData", encode64(out))

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

func menuLogueado(eleccion *int, username string) {
	menuLogueado :=
		`		
		[ 1 ] Gestionar contraseñas de sitios web
		[ 2 ] Gestionar tarjetas de cŕedito
		[ 3 ] Gestionar notas
		[ 4 ] Eliminar tu usuario
		[ 5 ] Cerrar sesión
		¿Qué prefieres?
	`
	fmt.Println()
	fmt.Print(fmt.Sprintf("Bienvenido %s.", username))
	fmt.Print(menuLogueado)
	fmt.Scanln(eleccion)
}

func logueado(client *http.Client, username string) {
	var eleccion int
	menuLogueado(&eleccion, username)

	for {
		switch eleccion {
		case 1: //GestionContraseñas
			gestionContraseñas(client, username)

		case 2: //GestionCard
			gestionTarjetas(client, username)

		case 3: //GestionNote
			gestionNotas(client, username)

		case 4:
			data := url.Values{} // estructura para contener los valores

			data.Set("cmd", "deleteUser") // comando (string)

			data.Set("username", username) // usuario (string)

			r, err := client.PostForm("https://localhost:8080", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

		case 5:
			fmt.Println("Hasta la vista.")
			return

		default:
			fmt.Println("No has seleccionado una opcion correcta.")
		}
	}
}

func menuGestionContraseña(eleccion *int) {
	menuGestionContraseña :=
		`		
		[ 1 ] Añadir contraseñas
		[ 2 ] Listar contraseñas
		[ 3 ] Modificar contraseñas
		[ 4 ] Eliminar contraseñas
		[ 5 ] Ir atrás
		¿Qué prefieres?
	`
	fmt.Print(menuGestionContraseña)
	fmt.Scanln(eleccion)
}

func gestionContraseñas(client *http.Client, username string) {
	var eleccion int
	menuGestionContraseña(&eleccion)

	for {
		switch eleccion {
		case 1: //addPassword
			data := url.Values{} // estructura para contener los valores

			data.Set("cmd", "addPassword") // comando (string)

			fmt.Print("Inserte URL: ")
			var url string
			fmt.Scanf("%s", &url)

			fmt.Print("Inserte usuario: ")
			var user string
			fmt.Scanf("%s", &user)

			fmt.Print("¿Quieres generar una contraseña aleatoria?(s/n) ")
			var opcion string
			fmt.Scanf("%s", &opcion)

			var contraseña string
			if opcion == "s" {
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
				contrasenyaa, err := password.Generate(long, numDigitos, numSimbolos, !upperLower, repeatCharacers)
				if err != nil {
					log.Fatal(err)
				}
				contraseña = contrasenyaa
			} else {
				fmt.Print("Introduce contraseña: ")
				fmt.Scanf("%s", &contraseña)
			}

			fmt.Printf("La contraseña generada es: ")
			fmt.Println(contraseña)

			data.Set("username", username) // usuario (string)
			data.Set("user", user)
			data.Set("site", url)
			data.Set("contraseña", contraseña)

			r, err := client.PostForm("https://localhost:8080", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return
		case 2: //List password

		case 3: //Modify password

		case 4: //Delete password

		case 5:
			return
		default:
			fmt.Println("No has seleccionado una opcion correcta.")

		}
	}
}

func menuGestionTarjetas(eleccion *int) {
	menuGestionTarjetas :=
		`		
		[ 1 ] Añadir tarjeta	
		[ 2 ] Listar tarjetas
		[ 3 ] Modificar tarjeta
		[ 4 ] Eliminar tarjeta
		[ 5 ] Ir atrás
		¿Qué prefieres?
	`
	fmt.Print(menuGestionTarjetas)
	fmt.Scanln(eleccion)
}

func gestionTarjetas(client *http.Client, username string) {
	var eleccion int
	menuGestionTarjetas(&eleccion)

	for {
		switch eleccion {
		case 1: //addCard
			data := url.Values{} // estructura para contener los valores

			data.Set("cmd", "addCreditCard") // comando (string)

			fmt.Print("Inserte número de la tarjeta: ")
			var pan string
			fmt.Scanf("%s", &pan)

			fmt.Print("Inserte CCV: ")
			var ccv string
			fmt.Scanf("%s", &ccv)

			fmt.Print("Inserte mes de caducidad: ")
			var month string
			fmt.Scanf("%s", &month)

			fmt.Print("Inserte año de caducidad: ")
			var year string
			fmt.Scanf("%s", &year)

			data.Set("username", username) // usuario (string)
			data.Set("pan", pan)
			data.Set("ccv", ccv)
			data.Set("month", month)
			data.Set("year", year)

			r, err := client.PostForm("https://localhost:8080", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 2: //List cards

		case 3: //Modify card

		case 4: //Delete card

		case 5:
			return

		default:
			fmt.Println("No has seleccionado una opcion correcta.")

		}
	}
}

func menuGestionNotas(eleccion *int) {
	menuGestionNotas :=
		`		
		[ 1 ] Añadir nota
		[ 2 ] Listar notas
		[ 3 ] Modificar nota		
		[ 4 ] Eliminar nota
		[ 5 ] Ir atrás
		¿Qué prefieres?
	`
	fmt.Print(menuGestionNotas)
	fmt.Scanln(eleccion)
}

func gestionNotas(client *http.Client, username string) {
	notas := make(map[int]notesData)

	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "Notes") // comando (string)
	data.Set("username", username)

	r, err := client.PostForm("https://localhost:8080/notes", data) // enviamos por POST
	chk(err)

	//--------- Con esto recojo del servidor las notas y las convierto al struct
	body, err := ioutil.ReadAll(r.Body)

	dec := json.NewDecoder(strings.NewReader(string(body)))

	for {
		var m resp
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(decode64(m.Msg), &notas) // Con esto paso al map notas lo que recojo en el servidor
	}
	//------------------------------------------------------------------------

	var eleccion int
	menuGestionNotas(&eleccion)

	for {
		switch eleccion {
		case 1: //add note
			newNota := notesData{}

			data := url.Values{} // estructura para contener los valores

			data.Set("cmd", "modifyNotes") // comando (string)

			fmt.Print("Inserte nota: ")
			var text string
			fmt.Scanf("%s", &text)

			fmt.Print("Inserte fecha: ")
			var date string
			fmt.Scanf("%s", &date)

			newNota.Text = text
			newNota.Date = date

			notas[len(notas)+1] = newNota
			out, err := json.Marshal(notas)
			if err != nil {
				panic(err)
			}

			data.Set("username", username)
			data.Set("notas", encode64(out))

			r, err := client.PostForm("https://localhost:8080/notes", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 2: //List notes
			for k, v := range notas {
				fmt.Println(k, "-. Texto: ", v.Text, ", fecha: ", v.Date)
			}

			return

		case 3: //Modify note
			modifyNota := notesData{}
			data := url.Values{} // estructura para contener los valores

			for k, v := range notas {
				fmt.Println(k, "-. Texto: ", v.Text, ", fecha: ", v.Date)
			}

			data.Set("cmd", "modifyNotes") // comando (string)

			fmt.Print("¿Que nota quieres editar?(num) ")
			var index int
			fmt.Scanf("%d", &index)
			delete(notas, index)

			fmt.Print("Inserte nueva nota: ")
			var text string
			fmt.Scanf("%s", &text)

			fmt.Print("Inserte nueva fecha: ")
			var date string
			fmt.Scanf("%s", &date)

			modifyNota.Text = text
			modifyNota.Date = date

			notas[index] = modifyNota
			out, err := json.Marshal(notas)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(out))

			data.Set("username", username)
			data.Set("notas", encode64(out))

			r, err := client.PostForm("https://localhost:8080/notes", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()
			return

		case 4: //Delete note
			data := url.Values{} // estructura para contener los valores

			for k, v := range notas {
				fmt.Println(k, "-. Texto: ", v.Text, ", fecha: ", v.Date)
			}

			data.Set("cmd", "modifyNotes") // comando (string)

			fmt.Print("¿Que nota quieres borrar?(num) ")
			var index int
			fmt.Scanf("%d", &index)
			delete(notas, index)

			out, err := json.Marshal(notas)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(out))

			data.Set("username", username)
			data.Set("notas", encode64(out))

			r, err := client.PostForm("https://localhost:8080/notes", data) // enviamos por POST
			chk(err)
			io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
			fmt.Println()

			return

		case 5:
			return

		default:
			fmt.Println("No has seleccionado una opcion correcta.")

		}
	}
}

/*
func main() {
	m := make(map[int]string)
	fmt.Println("Values in map (after creating): ", m)
	m[0] = "ABC"
	m[1] = "QR"
	m[2] = "XYZ"

	fmt.Println("Length of map: ", len(m))
	fmt.Println("Values in map(after adding values): ", m)

	m[1] = "LMN"
	fmt.Println("Values in map (after updating): ", m)

	delete(m, 1)
	fmt.Println("Values in map: ", m)

	for k, v := range m {
		fmt.Println("Key :", k, " Value :", v)
	}

	fmt.Println("Value for not existing key : ", m[3])
}
*/
