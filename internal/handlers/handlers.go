package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ava/notes-api/internal/config"
	"github.com/ava/notes-api/internal/models"
	"github.com/ava/notes-api/internal/storage"
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

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *Handler) HandleNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listNotes(w, r)
	case http.MethodPost:
		h.createNote(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.GetAll()
	if err != nil {
		log.Printf("Error fetching notes: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// basic validation
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) HandleNoteByID(w http.ResponseWriter, r *http.Request) {
	// extract ID from path
	// will fix later once i decide on how to handle ids
	path := strings.TrimPrefix(r.URL.Path, "/notes/")
	if path == "" {
		http.Error(w, "Note ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getNoteByID(w, r, path)
	case http.MethodDelete:
		// TODO: delete endpoint not implemented yet
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getNoteByID(w http.ResponseWriter, r *http.Request, id string) {
	note, err := h.store.GetByID(id)
	if err == storage.ErrNotFound {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error fetching note: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

