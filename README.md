# SDS

#### Pasos previos para poder arrancar la aplicación:
- Instalar [crypto](https://github.com/golang/crypto):
> go get -u golang.org/x/crypto/scrypt

- Instalar [mysql](https://github.com/go-sql-driver/mysql):
> go get -u github.com/go-sql-driver/mysql

- Instalar [go-password](https://github.com/sethvargo/go-password/password):
> go get -u github.com/sethvargo/go-password/password

- Instalar [aurora](https://github.com/logrusorgru/aurora):
> go get -u github.com/logrusorgru/aurora

#### Para arrancar el servidor, situarse en la carpeta de este y ejecutar:
> go run *.go [Puerto del servidor]

#### Para arrancar el cliente, situarse en la carpeta de este y ejecutar:
> go run *.go [IP del servidor] [Puerto del servidor]
