package routes

import (
	"encoding/json"
	"moscow_dep_task/db"
	"moscow_dep_task/types"
	"net/http"
	"sync/atomic"
	"time"
)

func SemaphoreAnalytics(w http.ResponseWriter, req *http.Request) {
	go func() {
		types.Log.Debugf("Semaphore analytics route start")
		types.C <- 1
		Analytics(w, req)
	}()
	<-types.C
}

func Analytics(w http.ResponseWriter, req *http.Request) {
	types.Log.Debugf("Analytics route start")
	atomic.AddInt32(&types.Counter, 1)
	go func() {
		user_id := req.Header.Get("X-Tantum-Authorization")
		data := types.Data{
			Header: req.Header,
			Body:   req.Body,
		}
		databytes, err := json.Marshal(data)
		if err != nil {
			types.Log.Errorf("Failed to marshal data into JSON: %s", err)
			return
		}
		types.Log.Debugf("Data marshalled into JSON")
		timestamp := time.Now()
		err = db.CreateRow(types.Conn, user_id, databytes, timestamp)
		if err != nil {
			types.Log.Errorf("Failed to create a user data row: %s", err)
			return
		}
		types.Log.Infof("User data row created")
	}()
	_, err := w.Write([]byte("Request acquired, response sent"))
	if err != nil {
		types.Log.Errorf("Failed to write a response to the user: %s", err)
		return
	}
	return
}
