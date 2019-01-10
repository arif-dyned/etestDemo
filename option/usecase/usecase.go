package usecase

import (
	"context"
	"time"

	"github.com/DynEd/etest/domain"
)

type optionUsecase struct {
	optionRepository domain.OptionRepository
	timeout            time.Duration
}

// NewOptionUsecase will create new an questionUsecase object representation of domain.OptionUsecase interface
func NewOptionUsecase(or domain.OptionRepository, timeout time.Duration) domain.OptionUsecase {
	return &optionUsecase{
		optionRepository: or,
		timeout:       timeout,
	}
}

// Fetch returns questions with it's related repository
// Note: when using inmem repo, the usecase is a bit of useless
func (u *optionUsecase) Fetch(ctx context.Context, num int64) ([]*domain.Option, error) {
	return u.optionRepository.Fetch(ctx, num)
}

// GetByID returns single question based on given ID
func (u *optionUsecase) GetByID(ctx context.Context, id int) (*domain.Option, error) {
	return u.optionRepository.GetByID(ctx, id)
}

// Store stores question into repository
func (u *optionUsecase) Store(question domain.Option) (int64, error) {
	return u.optionRepository.Store(question)
}

// Update updates existing question into new question
func (u *optionUsecase) Update(question domain.Option) (domain.Option, error) {
	return u.optionRepository.Update(question)
}

// Delete removes existing question from repository based on given ID
func (u *optionUsecase) Delete(id int) error {
	return u.optionRepository.Delete(id)
}
