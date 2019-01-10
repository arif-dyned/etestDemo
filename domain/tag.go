package domain

import (
	"time"
)

// Tag define the skill metatag of question
type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TagRepository represent the question's tag repository contract
type TagRepository interface {
	Fetch() ([]Tag, error)
	GetByID(id int) (Tag, error)
	Store(tag Tag) error
	Update(tag Tag) error
	Delete(id int) error
}
