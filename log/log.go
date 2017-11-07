package log

import (
	"fmt"
	"log"
	"os"
)

// Debug enabled debug logging
var Debug = false

var logger = log.New(os.Stderr, "", log.LstdFlags)

// Log produces log messages if the Debug flag is set to true
func Log(msg string, args ...interface{}) {

	if Debug {
		logger.Printf("[DEBUG] %s", fmt.Sprintf(msg, args...))
	}

}
