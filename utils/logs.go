package utils

import (
	"log"
	"os"
)

var (
	Warn  *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func Logs() {
	Warn = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
