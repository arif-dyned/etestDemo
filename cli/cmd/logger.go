package cmd

import (
	"fmt"
	"log"
	"os"
)

type logger struct {
	verbose bool
}

// Print returns log output
func (l logger) Printf(format string, v ...interface{}) {
	if l.verbose {
		log.Printf(format, v...)
	} else {
		fmt.Fprintf(os.Stderr, format, v...)
	}
}

// Verbose return logger verbose flag
func (l logger) Verbose() bool {
	return l.verbose
}
