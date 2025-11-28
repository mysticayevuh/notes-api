package storage

import (
	"errors"
	"sync"

	"github.com/ava/notes-api/internal/models"
)

var (
	ErrNotFound = errors.New("note not found")
)

// simple memory storage implementation
// TODO: convert to postgres with supabase
type MemoryStore struct {
	mu    sync.RWMutex
	notes map[string]*models.Note
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		notes: make(map[string]*models.Note),
	}
}

func (s *MemoryStore) Create(note *models.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.notes[note.ID] = note
	return nil
}

func (s *MemoryStore) GetByID(id string) (*models.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	note, exists := s.notes[id]
	if !exists {
		return nil, ErrNotFound
	}
	return note, nil
}

func (s *MemoryStore) GetAll() ([]*models.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// not proud of this allocation but it works for now
	result := make([]*models.Note, 0, len(s.notes))
	for _, note := range s.notes {
		result = append(result, note)
	}
	return result, nil
}

func (s *MemoryStore) Update(id string, note *models.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.notes[id]; !exists {
		return ErrNotFound
	}
	
	s.notes[id] = note
	return nil
}

// not yet decided on the return signature
func (s *MemoryStore) Delete(id string) error {
	// implement this
	return nil
}

