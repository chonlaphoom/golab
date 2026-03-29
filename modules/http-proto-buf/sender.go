package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

func sender(ctx context.Context) {
	log.Println("Starting sender...")

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			b := scanner.Bytes()
			res, err := http.Post("http://localhost:7777/message", "application/json", bytes.NewReader(b))
			fmt.Printf("Sent: %s\n", string(b))
			if err != nil {
				log.Printf("Error sending request: %v", err)
				continue
			}
			res.Body.Close()
			if res.StatusCode != http.StatusOK {
				log.Printf("Received non-OK response: %s", res.Status)
			}
			fmt.Printf("Received response: %s\n", res.Status)
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
		log.Println("Input closed, exiting sender...")
	}()

	<-ctx.Done()
	os.Stdin.Close()

	log.Println("Sender Exiting..")
}
