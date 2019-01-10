package usecase

import (
	"context"
	"time"

	"github.com/DynEd/etest/domain"
)

type questionUsecase struct {
	questionRepository domain.QuestionRepository
	timeout            time.Duration
}

// NewQuestionUsecase will create new an questionUsecase object representation of domain.QuestionUsecase interface
func NewQuestionUsecase(or domain.QuestionRepository, timeout time.Duration) domain.QuestionUsecase {
	return &questionUsecase{
		questionRepository: or,
		timeout:       timeout,
	}
}

// Fetch returns questions with it's related repository
// Note: when using inmem repo, the usecase is a bit of useless
func (u *questionUsecase) Fetch(ctx context.Context, num int64) ([]*domain.Question, error) {
	return u.questionRepository.Fetch(ctx, num)
}

// GetByID returns single question based on given ID
func (u *questionUsecase) GetByID(ctx context.Context, id int) (*domain.Question, error) {
	return u.questionRepository.GetByID(ctx, id)
}

// Store stores question into repository
func (u *questionUsecase) Store(question domain.Question) (int64, error) {
	return u.questionRepository.Store(question)
}

// Update updates existing question into new question
func (u *questionUsecase) Update(question domain.Question) (domain.Question, error) {
	return u.questionRepository.Update(question)
}

// Delete removes existing question from repository based on given ID
func (u *questionUsecase) Delete(id int) error {
	return u.questionRepository.Delete(id)
}
