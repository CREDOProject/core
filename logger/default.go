package logger

import (
	"log"
	"os"
)

// Global logger.
var logger *log.Logger

// Gets global logger.
func Get() *log.Logger {
	if logger != nil {
		return logger
	}

	logger = log.New(os.Stdout, "", log.Ltime)
	return logger
}
