package domain

import (
	"context"
	"github.com/mholt/binding"
	"net/http"
)

// Question represent Question's data model
type Question struct {
	ID               int64  `json:"id"`
	OptionMode       string `json:"option_mode,omitempty"`
	Instructions     string `json:"instructions,omitempty"`
	Comments         string `json:"comments,omitempty"`
	Type             string `json:"type,omitempty"`
	TagName          string `json:"tag_name,omitempty"`
	QuestionSlotType string `json:"question_slot_type,omitempty"`
	QuestionSlotFile string `json:"question_slot_file,omitempty"`
	QuestionSlotText string `json:"question_slot_text,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

// FieldMap binds HTTP request into Question
func (o *Question) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID:               "id",
		&o.OptionMode:       "option_mode",
		&o.Instructions:     "instructions",
		&o.Comments:         "comments",
		&o.Type:             "type",
		&o.TagName:          "tag_name",
		&o.QuestionSlotType: "question_slot_type",
		&o.QuestionSlotFile: "question_slot_file",
		&o.QuestionSlotText: "question_slot_text",
	}
}

// Validate perform data validation
func (o *Question) Validate(req *http.Request, errs binding.Errors) error {
	if o.Type == "" {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"message"},
			Classification: "ComplaintError",
			Message:        "Type can't be blank",
		})
	}

	return errs
}

// QuestionRepository represent the tags's repository contract
type QuestionRepository interface {
	Fetch(ctx context.Context, num int64) ([]*Question, error)
	GetByID(ctx context.Context, id int) (*Question, error)
	Store(Question Question) (int64, error)
	Update(Question Question) (Question, error)
	Delete(id int) error
}

// QuestionUsecase represent the tags's usecases
type QuestionUsecase interface {
	Fetch(ctx context.Context, num int64) ([]*Question, error)
	GetByID(ctx context.Context, id int) (*Question, error)
	Store(Question Question) (int64, error)
	Update(Question Question) (Question, error)
	Delete(id int) error
}
