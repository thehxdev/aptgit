package log

import (
	"log"
	"os"
)

var (
	Err *log.Logger = log.New(os.Stderr, "[ERROR] ", log.Lshortfile|log.Ldate|log.Ltime)
	Inf *log.Logger = log.New(os.Stdout, "[INFO] ", log.Lshortfile|log.Ldate|log.Ltime)
)
