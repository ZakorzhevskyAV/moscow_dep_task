package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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

	types.Log = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	LogFormatter := os.Getenv("LOG_FORMATTER")
	switch LogFormatter {
	case "json":
		types.Log.Formatter = new(logrus.JSONFormatter)
	case "text":
		types.Log.Formatter = new(logrus.TextFormatter)
	default:
		types.Log.Formatter = new(logrus.TextFormatter)
	}

	LogLevel := os.Getenv("LOG_LEVEL")
	switch LogLevel {
	case "fatal":
		types.Log.Level = logrus.FatalLevel
	case "error":
		types.Log.Level = logrus.ErrorLevel
	case "warn":
		types.Log.Level = logrus.WarnLevel
	case "debug":
		types.Log.Level = logrus.DebugLevel
	case "info":
		types.Log.Level = logrus.InfoLevel
	case "panic":
		types.Log.Level = logrus.PanicLevel
	case "trace":
		types.Log.Level = logrus.TraceLevel
	default:
		types.Log.Level = logrus.InfoLevel
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	types.Conn, err = sql.Open("postgres", connString)
	if err != nil {
		types.Log.Fatalf("Failed to connect to DB, panicking")
	}

	err = types.Conn.Ping()
	if err != nil {
		types.Log.Fatalf("Failed to check the connection to DB via ping, panicking")
	}
}

func main() {
	maxconn, err := strconv.Atoi(os.Getenv("MAX_CONNECTIONS"))
	if err != nil || maxconn == 0 {
		types.Log.Debugf("Registering the analytics route")
		http.HandleFunc("/analytics", routes.Analytics)
	} else {
		types.C = make(chan int, maxconn)
		types.Log.Debugf("Registering the semaphore analytics route")
		http.HandleFunc("/analytics", routes.SemaphoreAnalytics)
	}
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for tick := range ticker.C {
			val := atomic.SwapInt32(&types.Counter, 0)
			types.Log.Debugf("(%v) %d RPS", tick, val)
		}
	}()
	err = http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		log.Fatalf("Failed to listen and serve, panicking")
	}
}
