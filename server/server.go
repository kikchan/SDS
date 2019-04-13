package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
	Name     string
	Surname  string
	Email    string
	/* FALTA AÑADIR ESTO
	Hash []byte            // hash de la contraseña
	Salt []byte            // sal para la contraseña
	Data map[string]string // datos adicionales del usuario
	*/
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
func response(w io.Writer, ok bool, msg string) {
	r := resp{Ok: ok, Msg: msg}    // formateamos respuesta
	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante
}

var usersBD []user

func main() {
	//var port = "8080"

	//HACER SELECT DE LA BD E INTRODUCIR USUARIOS EN usersBD

	http.HandleFunc("/", handler) // asignamos un handler global

	// escuchamos el puerto 8080 con https y comprobamos el error
	// Para generar certificados autofirmados con openssl usar:
	//    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=ES/ST=Alicante/L=Alicante/O=UA/OU=Org/CN=www.ua.com"
	chk(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil))

	/*
		DALUsersTest()
		DALCardsTest()
		DALNotesTest()

		if len(os.Args) == 2 {
			port = os.Args[1]
			fmt.Println("Server awaiting connections from port: " + port)
		} else {
			fmt.Println("Server awaiting connections from port: " + port + " (default)")
		}

		//Server is in listening mode
		ln, err := net.Listen("tcp", "localhost:"+port)
		chk(err)

		defer ln.Close()

		//Infinite loop
		for {
			//Accept every single user request
			conn, err := ln.Accept()
			chk(err)

			//Launch a concurrent lambda function
			go func() {
				//Gets the user's port
				_, port, err := net.SplitHostPort(conn.RemoteAddr().String())
				chk(err)

				fmt.Println("Connection: ", conn.LocalAddr(), " <--> ", conn.RemoteAddr())

				scanner := bufio.NewScanner(conn)

				//Scans the connection and reads the message
				for scanner.Scan() {
					//Print the user's message
					fmt.Println("Client[", port, "]: ", scanner.Text())

					//Send "ACK" to client
					fmt.Fprintln(conn, "ack: ", scanner.Text())
				}

				conn.Close()
				fmt.Println("Closed[", port, "]")
			}()
		}
	*/
}

func handler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()                              // es necesario parsear el formulario
	w.Header().Set("Content-Type", "text/plain") // cabecera estándar

	switch req.Form.Get("cmd") { // comprobamos comando desde el cliente
	case "login": // ** login
		fmt.Println(req.Form.Get("cmd"))
		code, userJSON := findUser(req.Form.Get("user"))

		fmt.Println(code)
		fmt.Println(userJSON)

		//password := decode64(req.Form.Get("pass")) // obtenemos la contraseña
		//hash, _ := scrypt.Key(password, user.Salt, 16384, 8, 1, 32) // scrypt(contraseña)
		//if bytes.Compare(u.Hash, hash) != 0 { // comparamos
		//	response(w, false, "Credenciales inválidas")
		//	return
		//}

		response(w, true, "Credenciales válidas")

	case "register": // ** registro
		code, _ := findUser(req.Form.Get("user"))
		if code == -1 {
			u := user{}
			u.Username = req.Form.Get("user") // username
			u.Name = req.Form.Get("name")
			u.Surname = req.Form.Get("surname")
			u.Email = req.Form.Get("email")
			u.Password = req.Form.Get("pass")
			//u.Salt = make([]byte, 16)                  // sal (16 bytes == 128 bits)
			//rand.Read(u.Salt)                          // la sal es aleatoria
			//u.Data = make(map[string]string)           // reservamos mapa de datos de usuario
			//u.Data["private"] = req.Form.Get("prikey") // clave privada
			//u.Data["public"] = req.Form.Get("pubkey")  // clave pública
			//password := decode64(req.Form.Get("pass")) // contraseña (keyLogin)

			code, userJSON := createUser(u.Username, u.Password, u.Name, u.Surname, u.Email)

			fmt.Println(code)
			fmt.Println(userJSON)

			// "hasheamos" la contraseña con scrypt
			//u.Hash, _ = scrypt.Key(password, u.Salt, 16384, 8, 1, 32)

			response(w, true, "Usuario registrado")
		} else {
			response(w, false, "Usuario ya registrado")
		}
	default:
		response(w, false, "Comando inválido")
	}
}
