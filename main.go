// main.go

package main

func main() {
	a := App{}
	dbUname := "postgres"
	dbPass := "postgres"
	dbName := "postgres"

	a.Initialize(dbUname, dbPass, dbName)
	a.Run(":8010")
}
