package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func db() {
	db, err := sql.Open("mysql", "sds:sds@tcp(185.207.145.237:3306)/sds")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//insert, err := db.Query("INSERT INTO sds.users (username, password, name, surname) VALUES (sds, sds, SDS, SDS);")
	insert, err := db.Query("INSERT INTO users VALUES ('sds', 'sds', 'SDS', 'SDS');")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var puerto = "8080"

	if len(os.Args) == 2 {
		puerto = os.Args[1]
		fmt.Println("Servidor escuchando por el puerto: " + puerto)
	} else {
		fmt.Println("Servidor escuchando por el puerto: " + puerto + " (por defecto)")
	}

	ln, err := net.Listen("tcp", "localhost:"+puerto) // escucha en espera de conexión
	chk(err)
	defer ln.Close() // nos aseguramos que cerramos las conexiones aunque el programa falle

	for { // búcle infinito, se sale con ctrl+c
		conn, err := ln.Accept() // para cada nueva petición de conexión
		chk(err)
		go func() { // lanzamos un cierre (lambda, función anónima) en concurrencia

			_, port, err := net.SplitHostPort(conn.RemoteAddr().String()) // obtenemos el puerto remoto para identificar al cliente (decorativo)
			chk(err)

			fmt.Println("conexión: ", conn.LocalAddr(), " <--> ", conn.RemoteAddr())

			scanner := bufio.NewScanner(conn) // el scanner nos permite trabajar con la entrada línea a línea (por defecto)

			for scanner.Scan() { // escaneamos la conexión
				fmt.Println("cliente[", port, "]: ", scanner.Text()) // mostramos el mensaje del cliente
				fmt.Fprintln(conn, "ack: ", scanner.Text())          // enviamos ack al cliente
			}

			conn.Close() // cerramos al finalizar el cliente (EOF se envía con ctrl+d o ctrl+z según el sistema)
			fmt.Println("cierre[", port, "]")
		}()
	}
}
