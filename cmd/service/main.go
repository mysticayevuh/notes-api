package main

import (
	"log"
	"os"

	"github.com/ava/notes-api/internal/config"
	"github.com/ava/notes-api/internal/handlers"
	"github.com/ava/notes-api/internal/storage"
	"github.com/gofiber/fiber/v2"
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

	app := fiber.New(fiber.Config{
		// TODO: add more config options here
	})

	app.Get("/health", h.HealthCheck)
	app.Get("/notes", h.ListNotes)
	app.Post("/notes", h.CreateNote)
	app.Get("/notes/:id", h.GetNoteByID)
	app.Delete("/notes/:id", h.DeleteNote) // not fully completed yet

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

