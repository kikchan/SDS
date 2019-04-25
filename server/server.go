package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/scrypt"
)

//Database connection variables
var DB_IP string = "185.207.145.237"
var DB_Port string = "3306"
var DB_Protocol string = "tcp"
var DB_Name string = "sds2"
var DB_Username string = "sdsAppClient"
var DB_Password string = "qwerty123456"

//Error check function
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

type user struct {
	ID       int
	Username string
	Password string
	Hash     []byte // hash de la contraseña
	Salt     []byte // sal para la contraseña
	Data     string

	//Cambios a partir de aquí. Yo quitaría estos campos.
	Name    string
	Surname string
	Email   string
	DataOld map[string]string // datos adicionales del usuario
}

type notesData struct {
	Date string
	Text string
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

// respuesta del servidor
type resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
}

// función para escribir una respuesta del servidor
func response(w http.ResponseWriter, ok bool, msg string) {
	if ok == false {
		w.WriteHeader(404)
	}
	r := resp{Ok: ok, Msg: msg}    // formateamos respuesta
	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante
}

// función para escribir una respuesta del servidor
func response2(w http.ResponseWriter, msg string) {
	r := resp{Msg: msg}            // formateamos respuesta
	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante
}

func main() {
	/*
		if true {
			fmt.Println(createUser("kiril", "123456", "hash", "salt", "data"))
			fmt.Println(findUser("kiril"))
			fmt.Println(deleteUser("kiril"))
			fmt.Println(deleteUser("jose"))

			fmt.Println(createUser("kiril", "123456", "hash", "salt", "data"))

			fmt.Println(updatePassword("kiril", "contraseña"))
			fmt.Println(updateCard("kiril", "tarjeta"))
			fmt.Println(updateNote("kiril", "nota"))
		}
	*/

	http.HandleFunc("/", handler) // asignamos un handler global

	// escuchamos el puerto 8080 con https y comprobamos el error
	// Para generar certificados autofirmados con openssl usar:
	//    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=ES/ST=Alicante/L=Alicante/O=UA/OU=Org/CN=www.ua.com"
	chk(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()                              // es necesario parsear el formulario
	w.Header().Set("Content-Type", "text/plain") // cabecera estándar

	switch req.Form.Get("cmd") { // comprobamos comando desde el cliente
	case "login": // ** login
		_, _, user := findUser(req.Form.Get("user"))

		password := decode64(req.Form.Get("pass")) // obtenemos la contraseña

		hash, _ := scrypt.Key(password, user.Salt, 16384, 8, 1, 32) // scrypt(contraseña)
		if bytes.Compare(user.Hash, hash) != 0 {                    // comparamos
			response(w, false, "Credenciales inválidas")
			return
		}

		response(w, true, "Credenciales válidas")

	case "register": // ** registro
		code, _, _ := findUser(req.Form.Get("user"))
		if code == -3 {
			u := user{}
			u.Username = req.Form.Get("user") // username
			u.Data = req.Form.Get("userData")
			u.Password = req.Form.Get("pass")
			u.Salt = make([]byte, 16) // sal (16 bytes == 128 bits)
			rand.Read(u.Salt)         // la sal es aleatoria

			//----------------------------------------------------
			//		Aquí he cambiado u.Data por u.DataOld para que compile
			u.DataOld = make(map[string]string)           // reservamos mapa de datos de usuario
			u.DataOld["private"] = req.Form.Get("prikey") // clave privada
			u.DataOld["public"] = req.Form.Get("pubkey")  // clave pública
			password := decode64(req.Form.Get("pass"))    // contraseña (keyLogin)

			// "hasheamos" la contraseña con scrypt
			u.Hash, _ = scrypt.Key(password, u.Salt, 16384, 8, 1, 32)

			code, _ = createUser(u.Username, u.Password, encode64(u.Hash), encode64(u.Salt), u.Data)

			if code == 1 {
				response(w, true, "Usuario registrado")

			} else {
				response(w, true, "Usuario no se ha podido registrar")
			}
		} else {
			response(w, false, "Usuario ya registrado")
		}
	case "viewPassword": // ** View Password

	case "addPassword": // ** Add Password
		/*
			_, _, username := findUser(req.Form.Get("username"))

			user := req.Form.Get("user")
			site := req.Form.Get("site")
			contraseña := req.Form.Get("contraseña")
		*/

		//----------------------------------------------------
		//		Código viejo
		//
		//code, _ := createPassword(user, contraseña, username.Username, site)
		//
		//		Hay que meter los datos de la estructura en el campo data
		//----------------------------------------------------

		//if code == 1 {
		//	fmt.Println("Contraseña guardada con éxito.")
		//}

	case "addCreditCard": // ** Add credit card
		/*
			_, _, username := findUser(req.Form.Get("username"))

			pan := req.Form.Get("pan")
			ccv := req.Form.Get("ccv")
			monthString := req.Form.Get("month")
			yearString := req.Form.Get("year")

			month, _ := strconv.Atoi(monthString)
			year, _ := strconv.Atoi(yearString)
		*/

		//----------------------------------------------------
		//		Código viejo
		//
		//code, _ := createCard(pan, ccv, month, year, username.Username)
		//
		//		Hay que meter los datos de la estructura en el campo data
		//----------------------------------------------------

		//if code == 1 {
		//	fmt.Println("Tarjeta de crédito guardada con éxito.")
		//}

		return

	case "deleteUser":
		_, _, username := findUser(req.Form.Get("username"))

		code, _ := deleteUser(username.Username)

		if code == 1 {
			fmt.Println("Usuario eliminado con éxito.")
		} else {
			fmt.Println("El usuario no se puede eliminar.")
		}

		return

	case "Notes":

		a := &notesData{"21-05-1994", "hola"}

		out, err := json.Marshal(a)
		if err != nil {
			panic(err)
		}

		fmt.Println(a.Text)
		fmt.Println(string(out))
		response2(w, encode64(out))
		return

	default:
		response(w, false, "Comando inválido")
	}
}
