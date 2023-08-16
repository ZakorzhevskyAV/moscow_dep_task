package routes

import (
	"encoding/json"
	"log"
	"moscow_dep_task/db"
	"moscow_dep_task/types"
	"net/http"
	"time"
)

func SemaphoreAnalytics(w http.ResponseWriter, req *http.Request) {
	go func() {
		types.C <- 1
		Analytics(w, req)
		<-types.C
	}()
}

func Analytics(w http.ResponseWriter, req *http.Request) {
	go func() {
		user_id := req.Header.Get("X-Tantum-Authorization")
		data := types.Data{
			Header: req.Header,
			Body:   req.Body,
		}
		databytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to marshal data into JSON: %s\n", err)
			return
		}
		timestamp := time.Now()
		err = db.CreateRow(types.Conn, user_id, databytes, timestamp)
		if err != nil {
			log.Printf("Failed to create a user data row: %s\n", err)
			return
		}
	}()
	_, err := w.Write([]byte("aaa"))
	if err != nil {
		log.Printf("Failed to write a response to the user: %s\n", err)
		return
	}
}
