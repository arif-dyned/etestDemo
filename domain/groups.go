package domain

import (
	"context"
	"github.com/mholt/binding"
	"net/http"
)

// Groups represent Groups's data model
type Groups struct {
	ID           int64  `json:"id"`
	Name         string `json:"name,omitempty"`
	IsCore       string `json:"is_core,omitempty"`
	IsRefinement string `json:"is_refinement,omitempty"`
	Sequence     int    `json:"sequence,omitempty"`
	LevelId      string `json:"level_id,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// FieldMap binds HTTP request into Groups
func (o *Groups) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID:           "id",
		&o.Name:         "name",
		&o.IsCore:       "is_core",
		&o.IsRefinement: "is_refinement",
		&o.Sequence:     "sequence",
		&o.LevelId:      "level_id",
	}
}

// Validate perform data validation
func (o *Groups) Validate(req *http.Request, errs binding.Errors) error {
	if o.Name == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Name can't be blank",
		})
	} else if o.Sequence < 0 {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Sequence can't be blank",
		})
	} else if o.LevelId == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Level can't be blank",
		})
	}

	return errs
}

// GroupsRepository represent the groups's repository contract
type GroupsRepository interface {
	Fetch(ctx context.Context, num int64) ([]*Groups, error)
	GetByID(ctx context.Context, id int) (*Groups, error)
	Store(Groups Groups) (int64, error)
	Update(Groups Groups) (Groups, error)
	Delete(id int) error
}

// GroupsUsecase represent the groups's usecases
type GroupsUsecase interface {
	Fetch(ctx context.Context, num int64) ([]*Groups, error)
	GetByID(ctx context.Context, id int) (*Groups, error)
	Store(Groups Groups) (int64, error)
	Update(Groups Groups) (Groups, error)
	Delete(id int) error
}
