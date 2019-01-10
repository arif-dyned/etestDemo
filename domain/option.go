package domain

import (
	"context"
	"github.com/mholt/binding"
	"net/http"
)

// Option represent Option's data model
type Option struct {
	ID        int64  `json:"id"`
	Type      string `json:"type,omitempty"`
	Sequence  string `json:"sequence,omitempty"`
	IsCorrect string `json:"is_correct,omitempty"`
	IsLast    string `json:"is_last,omitempty"`
	File      string `json:"file,omitempty"`
	Text      string `json:"text,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// FieldMap binds HTTP request into Option
func (o *Option) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID:        "id",
		&o.Type:      "type",
		&o.Sequence:  "sequence",
		&o.IsCorrect: "is_correct",
		&o.IsLast:    "is_last",
		&o.File:      "file",
		&o.Text:      "text",
		&o.CreatedAt: "created_at",
		&o.UpdatedAt: "updated_at",
	}
}

// Validate perform data validation
func (o *Option) Validate(req *http.Request, errs binding.Errors) error {
	if o.Sequence == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Sequence can't be blank",
		})
	} else if o.IsCorrect == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "IsCorrect can't be blank",
		})
	} else if o.IsLast == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "IsLast can't be blank",
		})
	} else if o.Type == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Type can't be blank",
		})
	}

	return errs
}

// OptionRepository represent the option's repository contract
type OptionRepository interface {
	Fetch(ctx context.Context, num int64) ([]*Option, error)
	GetByID(ctx context.Context, id int) (*Option, error)
	Store(Option Option) (int64, error)
	Update(Option Option) (Option, error)
	Delete(id int) error
}

// OptionUsecase represent the option's usecases
type OptionUsecase interface {
	Fetch(ctx context.Context, num int64) ([]*Option, error)
	GetByID(ctx context.Context, id int) (*Option, error)
	Store(Option Option) (int64, error)
	Update(Option Option) (Option, error)
	Delete(id int) error
}
