// main.go

package main

import (
	"log"
)

func main() {
	a := App{}
	host := "database"
	port := 5432
	dbUname := "postgres"
	dbPass := "postgres"
	dbName := "postgres"

	serverPort := ":8010"

	a.Initialize(host, port, dbUname, dbPass, dbName)
	log.Print("Starting server on port ", serverPort, "\n")
	a.Run(serverPort)
}
