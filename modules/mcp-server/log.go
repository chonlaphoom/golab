package main

import (
	"log"
	"os"
)

func setupLog() {
	f, err := os.OpenFile("./mcp-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(f)
}
