package domain

import (
	"context"
	"github.com/mholt/binding"
	"net/http"
)

// Tags represent Tags's data model
type Tags struct {
	ID        int64  `json:"id"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// FieldMap binds HTTP request into Tags
func (o *Tags) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID:        "id",
		&o.Name:      "name",
	}
}

// Validate perform data validation
func (o *Tags) Validate(req *http.Request, errs binding.Errors) error {
	if o.Name == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Name can't be blank",
		})
	}

	return errs
}

// TagsRepository represent the tags's repository contract
type TagsRepository interface {
	Fetch(ctx context.Context, num int64) ([]*Tags, error)
	GetByID(ctx context.Context, id int) (*Tags, error)
	Store(Tags Tags) (int64, error)
	Update(Tags Tags) (Tags, error)
	Delete(id int) error
}

// TagsUsecase represent the tags's usecases
type TagsUsecase interface {
	Fetch(ctx context.Context, num int64) ([]*Tags, error)
	GetByID(ctx context.Context, id int) (*Tags, error)
	Store(Tags Tags) (int64, error)
	Update(Tags Tags) (Tags, error)
	Delete(id int) error
}
