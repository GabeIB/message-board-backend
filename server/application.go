// application.go

package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	a := App{}
	dbHost := os.Getenv("RDS_HOSTNAME")
	dbPort, err := strconv.Atoi(os.Getenv("RDS_PORT"))
	if err != nil{
		dbPort = 5432
	}
	dbName := os.Getenv("RDS_DB_NAME")
	dbUname := os.Getenv("RDS_UNAME")
	dbPass := os.Getenv("RDS_PASSWORD")

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "5000"
	}

	a.Initialize(dbHost, dbPort, dbUname, dbPass, dbName)
	log.Print("Starting server on port ", serverPort, "...\n")
	a.Run(serverPort)
}
