package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ava/notes-api/internal/config"
	"github.com/ava/notes-api/internal/handlers"
	"github.com/ava/notes-api/internal/storage"
)

func main() {
	// will move to config struct soon
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := config.New()

	store := storage.NewMemoryStore()
	
	// set up handlers
	h := handlers.New(store, cfg)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/notes", h.HandleNotes)
	mux.HandleFunc("/notes/", h.HandleNoteByID) // not fully completed yet
	
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// graceful shutdown
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down: %v", err)
	}
}

