package logger

import (
	"log"
	"os"
)

// Global logger.
var logger *log.Logger

// Get returns a global logger.
func Get() *log.Logger {
	if logger != nil {
		return logger
	}

	logger = log.New(os.Stdout, "", log.Ltime|log.LstdFlags|log.Lshortfile)
	return logger
}
