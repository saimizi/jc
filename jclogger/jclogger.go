package jclogger

import (
	"log"
	"os"
)

func NewErrLogger() *log.Logger {
	return log.New(os.Stderr, "ERROR: ", log.Lshortfile)
}

func NewWarnLogger() *log.Logger {
	return log.New(os.Stderr, "WARN: ", log.Lshortfile)
}

func NewInfoLogger() *log.Logger {
	return log.New(os.Stderr, "INFO: ", 0)
}

func NewDebugLogger() *log.Logger {
	return log.New(os.Stderr, "DEBUG: ", log.Lshortfile)
}
