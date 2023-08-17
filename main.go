package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"moscow_dep_task/routes"
	"moscow_dep_task/types"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

func init() {
	var err error

	//flag.Parse()
	//path := flag.String("config", "config.yaml", "Path to config file.")
	//err := config.ParseConfig(*path, types.Cfg)
	//
	//if err != nil {
	//	log.Printf("No config at the selected path or failed to parse config into struct; using default settings")
	//}
	//
	//types.Cfg = &types.Config{
	//	Server: types.Server{
	//		Host:           "",
	//		Port:           "8080",
	//		MaxConnections: 50,
	//	},
	//	Database: types.Database{
	//		User:     "pguser",
	//		Password: "pgpass",
	//		DB:       "pgdb",
	//		Host:     "postgres",
	//	},
	//}

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	types.Conn, err = sql.Open("postgres", connString)
	if err != nil {
		log.Panicf("Failed to connect to DB, panicking")
	}

	err = types.Conn.Ping()
	if err != nil {
		log.Panicf("Failed to check the connection to DB via ping, panicking")
	}
}

func main() {
	maxconn, err := strconv.Atoi(os.Getenv("MAX_CONNECTIONS"))
	if err != nil || maxconn == 0 {
		log.Printf("registering the analytics route")
		http.HandleFunc("/analytics", routes.Analytics)
	} else {
		types.C = make(chan int, maxconn)
		http.HandleFunc("/analytics", routes.SemaphoreAnalytics)
	}
	time.NewTicker(1 * time.Second)
	atomic.SwapInt32()
	err = http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		log.Panicf("Failed to listen and serve, panicking")
	}
}
