package db

import (
	"database/sql"
	"log"
	"time"
)

func CreateRow(conn *sql.DB, user_id string, data []byte, timestamp time.Time) error {
	_, err := conn.Query(`INSERT INTO UserData VALUES (?, ?, ?)`, user_id, data, timestamp)
	if err != nil {
		log.Printf("Failed to insert a data row: %s\n", err)
		return err
	}
	return err
}
