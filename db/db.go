package db

import (
	"database/sql"
	"moscow_dep_task/types"
	"time"
)

func CreateRow(conn *sql.DB, user_id string, data []byte, timestamp time.Time) error {
	_, err := conn.Exec("INSERT INTO UserData VALUES ($1, $2, $3);", user_id, data, timestamp)
	if err != nil {
		types.Log.Errorf("Failed to insert a data row: %s", err)
		return err
	}
	return err
}
