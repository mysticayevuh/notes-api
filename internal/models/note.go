package models

import "time"

// each note represents a single entry
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// expected request when creating a note
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

