package jc

import (
	"log"
	"os"
)

// NullLogger : null logger handler
func NullLogger() *log.Logger {
	f, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		log.Fatal("Faile to open /dev/null")
	}

	return log.New(f, "", 0)
}

// NewErrLogger  : create new error logger handler
func NewErrLogger() *log.Logger {
	return log.New(os.Stderr, "ERROR: ", 0)
}

// NewWarnLogger : create new warn logger handler
func NewWarnLogger() *log.Logger {
	switch os.Getenv("JCDBG") {
	case "error":
		break
	case "warn":
		return log.New(os.Stderr, "WARN: ", 0)
	case "info":
		return log.New(os.Stderr, "WARN: ", 0)
	case "debug":
		return log.New(os.Stderr, "WARN: ", 0)
	default:
		return log.New(os.Stderr, "WARN: ", 0)
	}

	return NullLogger()
}

// NewInfoLogger : create new info logger
func NewInfoLogger() *log.Logger {
	switch os.Getenv("JCDBG") {
	case "error":
		break
	case "warn":
		break
	case "info":
		return log.New(os.Stderr, "INFO: ", 0)
	case "debug":
		return log.New(os.Stderr, "INFO: ", 0)
	default:
		return log.New(os.Stderr, "INFO: ", 0)
	}

	return NullLogger()
}

// NewDebugLogger : new debug logger
func NewDebugLogger() *log.Logger {

	switch os.Getenv("JCDBG") {
	case "debug":
		return log.New(os.Stderr, "DEBUG: ", log.Lshortfile)
	default:
		break
	}

	return NullLogger()
}
