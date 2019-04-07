package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

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

func main() {
	var puerto = "8080"

	fmt.Println(createUser("kiril", "123456", "Kiril", "Gaydarov", "kvg1@alu.ua.es"))
	fmt.Println(findUser("kiril"))
	fmt.Println(updateUser("kiril", "654321", "kiril_gaydarov@gmail.com"))
	fmt.Println(findUser("kiril"))
	fmt.Println(deleteUser("kiril"))

	if len(os.Args) == 2 {
		puerto = os.Args[1]
		fmt.Println("Servidor escuchando por el puerto: " + puerto)
	} else {
		fmt.Println("Servidor escuchando por el puerto: " + puerto + " (por defecto)")
	}

	ln, err := net.Listen("tcp", "localhost:"+puerto) //escucha en espera de conexión
	chk(err)
	defer ln.Close() //nos aseguramos que cerramos las conexiones aunque el programa falle

	for { //búcle infinito, se sale con ctrl+c
		conn, err := ln.Accept() //para cada nueva petición de conexión
		chk(err)
		go func() { //lanzamos un cierre (lambda, función anónima) en concurrencia

			_, port, err := net.SplitHostPort(conn.RemoteAddr().String()) //obtenemos el puerto remoto para identificar al cliente (decorativo)
			chk(err)

			fmt.Println("conexión: ", conn.LocalAddr(), " <--> ", conn.RemoteAddr())

			scanner := bufio.NewScanner(conn) //el scanner nos permite trabajar con la entrada línea a línea (por defecto)

			for scanner.Scan() { //escaneamos la conexión
				fmt.Println("cliente[", port, "]: ", scanner.Text()) //mostramos el mensaje del cliente
				fmt.Fprintln(conn, "ack: ", scanner.Text())          //enviamos ack al cliente
			}

			conn.Close() //cerramos al finalizar el cliente (EOF se envía con ctrl+d o ctrl+z según el sistema)
			fmt.Println("cierre[", port, "]")
		}()
	}
}
