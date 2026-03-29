package main

import (
	"context"
	"errors"
	"http-proto-buf/generated/message_pb"
	"io"
	"log"
	"net/http"
	"time"

	proto "google.golang.org/protobuf/proto"
)

const MAX_MESSAGE_SIZE = 1048576 // 1MB

func receiver(ctx context.Context) {
	log.Println("Starting receiver...")
	port := "7777"
	address := ":" + port

	srv := &http.Server{
		Addr:    address,
		Handler: http.HandlerFunc(handler),
	}

	go func() {
		log.Println("Server starting on", address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Reciever Exiting...")
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/message" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, MAX_MESSAGE_SIZE)
	messageData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	message := &message_pb.Message{}
	err = proto.Unmarshal(messageData, message)
	if err != nil {
		http.Error(w, "Error unmarshaling body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received request: %s %s, %s", r.Method, r.URL.Path, message.Content)
	w.WriteHeader(http.StatusOK)
}
