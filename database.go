//contains functions for database management

package main

import (
	"log"
	"database/sql"
	"fmt"
)

//ensureTableExists creates the messages table in the database if there isn't already one.
func ensureTableExists(db *sql.DB) {
	if _, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		log.Fatal(err)
	}

	const tableCreationQuery = `CREATE TABLE IF NOT EXISTS messages
	(
	    id uuid PRIMARY KEY default uuid_generate_v1(),
	    name TEXT NOT NULL,
	    email TEXT NOT NULL,
	    text TEXT NOT NULL,
	    creation_time TIMESTAMPTZ NOT NULL default NOW()
	)`

	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS timestamp_desc_index ON messages (creation_time DESC)`); err != nil {
		log.Fatal(err)
	}
}

func clearTable(db *sql.DB) {
    db.Exec("DELETE FROM messages")
}

//loadDataFromCSV loads the specified filepath into the database.
func loadDataFromCSV(fileName string, db *sql.DB) error {
	query := fmt.Sprintf("COPY messages from '%s' DELIMITERS ',' CSV HEADER;", fileName)
	_, err := db.Exec(query)
	return err
	/*
	csvfile, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Could not open csv file")
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	//throw away first line of csv
	if _, err := r.Read(); err != nil{
		log.Fatal(err)
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		id := record[0]
		name := record[1]
		email := record[2]
		text := record[3]
		creation_time := record[4]
		query := fmt.Sprintf("INSERT INTO messages(id, name, email, text, creation_time) VALUES('%s', '%s', '%s', '%s', '%s');", id, name, email, text, creation_time)
		fmt.Println(query)
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil*/
}
