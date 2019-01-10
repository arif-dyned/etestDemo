package domain

import (
	"context"
	"github.com/mholt/binding"
	"net/http"
)

// QuestionSlot represent QuestionSlot's data model
type QuestionSlot struct {
	ID         int64  `json:"id"`
	Type       string `json:"type,omitempty"`
	File       string `json:"file,omitempty"`
	Text       string `json:"text,omitempty"`
	QuestionId string `json:"question_id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

// FieldMap binds HTTP request into QuestionSlot
func (o *QuestionSlot) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID:         "id",
		&o.Type:       "type",
		&o.File:       "file",
		&o.Text:       "text",
		&o.QuestionId: "question_id",
	}
}

// Validate perform data validation
func (o *QuestionSlot) Validate(req *http.Request, errs binding.Errors) error {
	if o.Type == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Type can't be blank",
		})
	}

	return errs
}

// QuestionSlotRepository represent the questionSlot's repository contract
type QuestionSlotRepository interface {
	Fetch(ctx context.Context, num int64) ([]*QuestionSlot, error)
	GetByID(ctx context.Context, id int, slotId int) (*QuestionSlot, error)
	Store(QuestionSlot QuestionSlot) (int64, error)
	Update(QuestionSlot QuestionSlot) (QuestionSlot, error)
	Delete(id int) error
}

// QuestionSlotUsecase represent the questionSlot's usecases
type QuestionSlotUsecase interface {
	Fetch(ctx context.Context, num int64) ([]*QuestionSlot, error)
	GetByID(ctx context.Context, id int, slotId int) (*QuestionSlot, error)
	Store(QuestionSlot QuestionSlot) (int64, error)
	Update(QuestionSlot QuestionSlot) (QuestionSlot, error)
	Delete(id int) error
}
