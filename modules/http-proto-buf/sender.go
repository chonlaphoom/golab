package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected 'sender' or 'receiver' as first argument")
		os.Exit(1)
	}
	first := os.Args[1]
	if first != "sender" && first != "receiver" {
		log.Fatalf("Expected 'sender' or 'receiver' as argument, got '%s'", first)
		os.Exit(1)
	}

	log.Println("Starting sender...")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("echo: %s\n", line)
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
		log.Println("Input closed, exiting sender...")
	}()

	<-ctx.Done()
	os.Stdin.Close()

	log.Println("Exiting..")
}
