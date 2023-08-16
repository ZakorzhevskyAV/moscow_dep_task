package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"moscow_dep_task/config"
	"moscow_dep_task/routes"
	"moscow_dep_task/types"
	"net/http"
)

func init() {
	flag.Parse()
	path := flag.String("config", "config.yaml", "Path to config file.")
	err := config.ParseConfig(*path, types.Cfg)

	if err != nil {
		log.Printf("No config at the selected path or failed to parse config into struct; using default settings")
	}

	types.Cfg = &types.Config{
		Server: types.Server{
			Host:           "",
			Port:           "8080",
			MaxConnections: 50,
		},
		Database: types.Database{
			User:     "pguser",
			Password: "pgpass",
			DB:       "pgdb",
			Host:     "postgres",
		},
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		types.Cfg.Database.User,
		types.Cfg.Database.Password,
		types.Cfg.Database.Host,
		types.Cfg.Database.DB)

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
	switch types.Cfg.Server.ConnLimit {
	case true:
		types.C = make(chan int, types.Cfg.Server.MaxConnections)
		http.HandleFunc("/analytics", routes.SemaphoreAnalytics)
	case false:
		http.HandleFunc("/analytics", routes.Analytics)
	default:
		http.HandleFunc("/analytics", routes.Analytics)
	}
	err := http.ListenAndServe(types.Cfg.Server.Host+":"+types.Cfg.Server.Port, nil)
	if err != nil {
		log.Panicf("Failed to listen and serve, panicking")
	}
}
