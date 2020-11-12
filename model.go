// model.go handles interaction with the database

package main

import (
	"fmt"
    "database/sql"
    "time"
)

//message is a structure that holds the information relating to a single post on the message-board.
type message struct {
    ID    string     `json:"id"`
    Name  string  `json:"name"`
    Email string `json:"email"`
    Text string `json:"text"`
    TimeStamp time.Time `json:"creation_time"`
}

func convertTime(timeString string) (time.Time, error) {
	layout := "2006-1-2T15:04:05-07:00" //layout as defined by sample CSV file
	t, err := time.Parse(layout, timeString)
	return t, err
}

//getMessage looks up a message by ID in the database and returns it.
func (m *message) getMessage(db *sql.DB) error {
	return db.QueryRow("SELECT name, email, text, creation_time FROM messages WHERE id=$1",
	        m.ID).Scan(&m.Name, &m.Email, &m.Text, &m.TimeStamp)
}

//updateMessage looks up a message by ID in the database and modifies the database fields to match the arguments.
func (m *message) updateMessage(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE messages SET name=$1, email=$2, text=$3 WHERE id=$4",
			m.Name, m.Email, m.Text, m.ID)

	return err
}

//loadDataFromCSV loads the specified filepath into the database.
//I could have done this in SQL with COPY FROM, but wanted to have more control over the inputs.
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

//createMessage adds a message to the database.
func (m *message) createMessage(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO messages(name, email, text) Values($1, $2, $3) RETURNING id, creation_time AS idstamp",
		m.Name, m.Email, m.Text).Scan(&m.ID, &m.TimeStamp)
	if err != nil {
		return err
	}
	return nil
}

//getMessages performs an sql query to return messages between $start and $start + $count not inclusive.
//Ordering of results is in reverse-chronological order.
func getMessages(db *sql.DB) ([]message, error) {
	rows, err := db.Query("SELECT id, name, email, text, creation_time FROM messages ORDER BY creation_time DESC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []message{}

	for rows.Next() {
		var m message
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Text, &m.TimeStamp); err != nil {
			return nil, err
		}
		messages = append(messages,m)
	}
	return messages, nil
}
