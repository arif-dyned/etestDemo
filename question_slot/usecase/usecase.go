package usecase

import (
	"context"
	"time"

	"github.com/DynEd/etest/domain"
)

type questionSlotUsecase struct {
	questionSlotRepository domain.QuestionSlotRepository
	timeout            time.Duration
}

// NewQuestionSlotUsecase will create new an questionUsecase object representation of domain.QuestionSlotUsecase interface
func NewQuestionSlotUsecase(or domain.QuestionSlotRepository, timeout time.Duration) domain.QuestionSlotUsecase {
	return &questionSlotUsecase{
		questionSlotRepository: or,
		timeout:       timeout,
	}
}

// Fetch returns questions with it's related repository
// Note: when using inmem repo, the usecase is a bit of useless
func (u *questionSlotUsecase) Fetch(ctx context.Context, num int64) ([]*domain.QuestionSlot, error) {
	return u.questionSlotRepository.Fetch(ctx, num)
}

// GetByID returns single question based on given ID
func (u *questionSlotUsecase) GetByID(ctx context.Context, id, slotId int) (*domain.QuestionSlot, error) {
	return u.questionSlotRepository.GetByID(ctx, id, slotId)
}

// Store stores question into repository
func (u *questionSlotUsecase) Store(question domain.QuestionSlot) (int64, error) {
	return u.questionSlotRepository.Store(question)
}

// Update updates existing question into new question
func (u *questionSlotUsecase) Update(question domain.QuestionSlot) (domain.QuestionSlot, error) {
	return u.questionSlotRepository.Update(question)
}

// Delete removes existing question from repository based on given ID
func (u *questionSlotUsecase) Delete(id int) error {
	return u.questionSlotRepository.Delete(id)
}
