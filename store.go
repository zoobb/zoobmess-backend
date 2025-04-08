package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Message struct {
	ID        int
	Username  string
	Message   string
	Timestamp int64 // тоже заменить на int64
}

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "chat.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		message TEXT NOT NULL,
		timestamp TEXT NOT NULL
	);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
}

func SaveMessage(msg Message) error {
	insertQuery := `INSERT INTO messages (username, message, timestamp) VALUES (?, ?, ?)`
	_, err := db.Exec(insertQuery, msg.Username, msg.Message, msg.Timestamp)
	return err
}

func LoadLastMessages(limit int) ([]Message, error) {
	selectQuery := `SELECT username, message, timestamp FROM messages ORDER BY id DESC LIMIT ?`
	rows, err := db.Query(selectQuery, limit)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Username, &msg.Message, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append([]Message{msg}, messages...)
	}
	return messages, nil
}
