package main

import (
	"log"
	"os"
)

const (
	VERSION = "1.0.0 \u00a9Inventos, Orel (RU), 2016"
)

// Entry point
func main() {
	// tell log to write to stdout, not stderr
	log.SetOutput(os.Stdout)

	if process_args() {
		os.Exit(0)
	}

	// run master process if no args specified
	run_master_process()
}
