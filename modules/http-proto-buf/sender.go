package main

import (
	"bufio"
	"bytes"
	"context"
	"http-proto-buf/generated/message_pb"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
)

func sender(ctx context.Context) {
	log.Println("Starting sender...")

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		url := "http://localhost:7777/message"
		contentType := "application/x-proto-content"

		for scanner.Scan() {
			b := scanner.Bytes()
			message := &message_pb.Message{
				Id:        "0",
				Content:   string(b),
				Timestamp: time.Now().Unix(),
			}
			data, err := proto.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
			if err != nil {
				log.Printf("Error creating request: %v", err)
				continue
			}
			req.Header.Set("Content-Type", contentType)
			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				log.Printf("Error sending request: %v", err)
				continue
			}

			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("Received non-OK response: %s", resp.Status)
			} else {
				log.Printf("Message sent successfully: %s", string(b))
			}
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
