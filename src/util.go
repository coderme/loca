package main

import (
	"fmt"
	"log"
	"os"
)

// exit writes error to stderr and exit
func exit(code int, s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
	os.Exit(code)
}

// exitf writes formated error msgs to stderr and exit
func exitf(code int, format string, s ...interface{}) {
	fmt.Fprintf(os.Stderr, format, s...)
	os.Exit(code)
}

// logf logs error and exit
func logf(code int, format string, s ...interface{}) {
	log.Printf(format, s...)
	os.Exit(code)
}
