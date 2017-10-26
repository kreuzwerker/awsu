package log

import (
	"fmt"
	"log"
)

var Debug = false

func Log(msg string, args ...interface{}) {

	if Debug {
		log.Printf("[DEBUG] %s", fmt.Sprintf(msg, args...))
	}

}
