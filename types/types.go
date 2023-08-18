package types

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Data struct {
	Header http.Header   `json:"header"`
	Body   io.ReadCloser `json:"body"`
}

var (
	Counter int32
	Conn    *sql.DB
	C       chan int
	Log     *logrus.Logger
)
