package domain

import (
	"context"
	"github.com/mholt/binding"
	"net/http"
)

// Types represent Types's data model
type Types struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// FieldMap binds HTTP request into Types
func (o *Types) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID:        "id",
		&o.Name:      "name",
	}
}

// Validate perform data validation
func (o *Types) Validate(req *http.Request, errs binding.Errors) error {
	if o.Name == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Name can't be blank",
		})
	}

	return errs
}

// TypesRepository represent the types's repository contract
type TypesRepository interface {
	Fetch(ctx context.Context, num int64) ([]*Types, error)
	GetByID(ctx context.Context, id int) (*Types, error)
	Store(Types Types) (int64, error)
	Update(Types Types) (Types, error)
	Delete(id int) error
}

// TypesUsecase represent the types's usecases
type TypesUsecase interface {
	Fetch(ctx context.Context, num int64) ([]*Types, error)
	GetByID(ctx context.Context, id int) (*Types, error)
	Store(Types Types) (int64, error)
	Update(Types Types) (Types, error)
	Delete(id int) error
}
