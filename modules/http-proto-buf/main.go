package main

import (
	"context"
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
	argument := os.Args[1]
	if argument != "sender" && argument != "receiver" {
		log.Fatalf("Expected 'sender' or 'receiver' as argument, got '%s'", argument)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	if argument == "sender" {
		sender(ctx)
	} else {
		receiver(ctx)
	}

	<-ctx.Done()
	os.Stdin.Close()

	log.Println("Exiting..")
}
