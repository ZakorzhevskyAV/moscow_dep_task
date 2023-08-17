package types

import (
	"database/sql"
	"io"
	"net/http"
)

//type Config struct {
//	Server   Server   `yaml:"address"`
//	Database Database `yaml:"database"`
//}
//
//type Server struct {
//	Host           string `yaml:"host"`
//	Port           string `yaml:"port"`
//	MaxConnections int    `yaml:"max_connections"`
//}
//
//type Database struct {
//	User     string `yaml:"user"`
//	Password string `yaml:"password"`
//	DB       string `yaml:"db"`
//	Host     string `yaml:"host"`
//}

type Data struct {
	Header http.Header   `json:"header"`
	Body   io.ReadCloser `json:"body"`
}

var (
	//Cfg  *Config
	Conn *sql.DB
	C    chan int
)
