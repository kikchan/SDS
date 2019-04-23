package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/scrypt"
)

//Database connection variables
var DB_IP string = "185.207.145.237"
var DB_Port string = "3306"
var DB_Protocol string = "tcp"
var DB_Name string = "sds"
var DB_Username string = "sds"
var DB_Password string = "sds"

//Error check function
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

type user struct {
	Username string
	Password string
	Hash     []byte // hash de la contraseña
	Salt     []byte // sal para la contraseña
	Name     string
	Surname  string
	Email    string
	Data     map[string]string // datos adicionales del usuario
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

var usersBD []user

func main() {
	//var port = "8080"

	/*******************************************************************************
	********* Los comentarios en inglés, por favor *********************************
	********************************************************************************
	***************Comenta las líneas de abajo para que no te tarde la *************
	********* vida ejecutando las pruebas ******************************************
	********************************************************************************/

	//DALUsersTest()
	//DALCardsTest()
	//DALNotesTest()
	//DALPasswordsTest()

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
		if code == -1 {
			u := user{}
			u.Username = req.Form.Get("user") // username
			u.Name = req.Form.Get("name")
			u.Surname = req.Form.Get("surname")
			u.Email = req.Form.Get("email")
			u.Password = req.Form.Get("pass")
			u.Salt = make([]byte, 16)                  // sal (16 bytes == 128 bits)
			rand.Read(u.Salt)                          // la sal es aleatoria
			u.Data = make(map[string]string)           // reservamos mapa de datos de usuario
			u.Data["private"] = req.Form.Get("prikey") // clave privada
			u.Data["public"] = req.Form.Get("pubkey")  // clave pública
			password := decode64(req.Form.Get("pass")) // contraseña (keyLogin)

			// "hasheamos" la contraseña con scrypt
			u.Hash, _ = scrypt.Key(password, u.Salt, 16384, 8, 1, 32)

			code, _ := createUser(u.Username, u.Password, encode64(u.Hash), encode64(u.Salt), u.Name, u.Surname, u.Email)

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
		_, _, username := findUser(req.Form.Get("username"))

		user := req.Form.Get("user")
		site := req.Form.Get("site")
		contraseña := req.Form.Get("contraseña")

		code, _ := createPassword(user, contraseña, username.Username, site)

		if code == 1 {
			fmt.Println("Contraseña guardada con éxito.")
		}

	case "addCreditCard": // ** Add credit card
		_, _, username := findUser(req.Form.Get("username"))

		pan := req.Form.Get("pan")
		ccv := req.Form.Get("ccv")
		monthString := req.Form.Get("month")
		yearString := req.Form.Get("year")

		month, _ := strconv.Atoi(monthString)
		year, _ := strconv.Atoi(yearString)

		code, _ := createCard(pan, ccv, month, year, username.Username)

		if code == 1 {
			fmt.Println("Tarjeta de crédito guardada con éxito.")
		}

		return
	default:
		response(w, false, "Comando inválido")
	}
}
