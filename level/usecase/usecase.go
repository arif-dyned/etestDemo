package usecase

import (
	"context"
	"time"

	"github.com/DynEd/etest/domain"
)

type levelUsecase struct {
	levelRepository domain.LevelRepository
	timeout            time.Duration
}

// NewLevelUsecase will create new an questionUsecase object representation of domain.LevelUsecase interface
func NewLevelUsecase(or domain.LevelRepository, timeout time.Duration) domain.LevelUsecase {
	return &levelUsecase{
		levelRepository: or,
		timeout:       timeout,
	}
}

// Fetch returns questions with it's related repository
// Note: when using inmem repo, the usecase is a bit of useless
func (u *levelUsecase) Fetch(ctx context.Context, num int64) ([]*domain.Level, error) {
	return u.levelRepository.Fetch(ctx, num)
}

// GetByID returns single question based on given ID
func (u *levelUsecase) GetByID(ctx context.Context, id int) (*domain.Level, error) {
	return u.levelRepository.GetByID(ctx, id)
}

// Store stores question into repository
func (u *levelUsecase) Store(question domain.Level) (int64, error) {
	return u.levelRepository.Store(question)
}

// Update updates existing question into new question
func (u *levelUsecase) Update(question domain.Level) (domain.Level, error) {
	return u.levelRepository.Update(question)
}

// Delete removes existing question from repository based on given ID
func (u *levelUsecase) Delete(id int) error {
	return u.levelRepository.Delete(id)
}
