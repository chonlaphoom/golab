package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

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
	message := make([]byte, r.ContentLength)
	r.Body.Read(message)
	defer r.Body.Close()

	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received message: " + string(message)))
}
