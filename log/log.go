package log

import (
	"fmt"
	"log"
	"os"
)

// Verbose enabled debug logging
var Verbose = false

var logger = log.New(os.Stderr, "", log.LstdFlags)

// Debug produces log messages if the Verbose flag is set to true
func Debug(msg string, args ...interface{}) {

	if Verbose {
		logger.Printf("[DEBUG] %s", fmt.Sprintf(msg, args...))
	}

}

// Info produces log messages regardless of Verbose
func Info(msg string, args ...interface{}) {
	logger.Printf(msg, args...)
}
