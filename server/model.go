// model.go holds message struct and message methods.

package main

import (
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
