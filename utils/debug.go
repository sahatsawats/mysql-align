package utils

import "fmt"

// global default variable
var debug bool = false

// init function to set debug. default is false.
func SetDebug(d bool) {
	debug = d
}


// main call function to print debug message. If debug not set to "true" -> do noting.
func Debug(messages string) {
	if debug {
		fmt.Printf("[DEBUG_MESSAGE]: %s\n", messages)
	}
}
