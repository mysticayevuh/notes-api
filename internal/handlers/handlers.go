package handlers

import (
	"log"
	"time"

	"github.com/ava/notes-api/internal/config"
	"github.com/ava/notes-api/internal/models"
	"github.com/ava/notes-api/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	store  Store
	config *config.Config
}

// will define more store types as needed
type Store interface {
	Create(note *models.Note) error
	GetByID(id string) (*models.Note, error)
	GetAll() ([]*models.Note, error)
	Update(id string, note *models.Note) error
	Delete(id string) error
}

func New(store Store, cfg *config.Config) *Handler {
	return &Handler{
		store:  store,
		config: cfg,
	}
}

func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (h *Handler) ListNotes(c *fiber.Ctx) error {
	notes, err := h.store.GetAll()
	if err != nil {
		log.Printf("Error fetching notes: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(notes)
}

func (h *Handler) CreateNote(c *fiber.Ctx) error {
	var req models.CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// basic validation
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	now := time.Now()
	note := &models.Note{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.store.Create(note); err != nil {
		log.Printf("Error creating note: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(note)
}

func (h *Handler) GetNoteByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Note ID required",
		})
	}

	note, err := h.store.GetByID(id)
	if err == storage.ErrNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Note not found",
		})
	}
	if err != nil {
		log.Printf("Error fetching note: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(note)
}

func (h *Handler) DeleteNote(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Note ID required",
		})
	}

	err := h.store.Delete(id)
	if err != nil {
		log.Printf("Error deleting note: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Note deleted successfully",
	})

	// return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
	// 	"error": "Not implemented",
	// })
}

