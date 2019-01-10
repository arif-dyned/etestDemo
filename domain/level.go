package domain

import (
	"context"
	"github.com/mholt/binding"
	"net/http"
)

// Level represent Level's data model
type Level struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	Sequence  int    `json:"sequence,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// FieldMap binds HTTP request into Level
func (o *Level) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID:        "id",
		&o.Name:      "name",
		&o.Sequence:  "sequence",
		&o.CreatedAt: "created_at",
		&o.UpdatedAt: "updated_at",
	}
}

// Validate perform data validation
func (o *Level) Validate(req *http.Request, errs binding.Errors) error {
	if o.Sequence == 0 {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Sequence can't be blank",
		})
	} else if o.Name == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "IsCorrect can't be blank",
		})
	}

	return errs
}

// LevelRepository represent the level's repository contract
type LevelRepository interface {
	Fetch(ctx context.Context, num int64) ([]*Level, error)
	GetByID(ctx context.Context, id int) (*Level, error)
	Store(Level Level) (int64, error)
	Update(Level Level) (Level, error)
	Delete(id int) error
}

// LevelUsecase represent the level's usecases
type LevelUsecase interface {
	Fetch(ctx context.Context, num int64) ([]*Level, error)
	GetByID(ctx context.Context, id int) (*Level, error)
	Store(Level Level) (int64, error)
	Update(Level Level) (Level, error)
	Delete(id int) error
}
