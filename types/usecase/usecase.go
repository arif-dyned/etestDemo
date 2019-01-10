package usecase

import (
	"context"
	"time"

	"github.com/DynEd/etest/domain"
)

type typesUsecase struct {
	typesRepository domain.TypesRepository
	timeout            time.Duration
}

// NewTypesUsecase will create new an questionUsecase object representation of domain.TypesUsecase interface
func NewTypesUsecase(or domain.TypesRepository, timeout time.Duration) domain.TypesUsecase {
	return &typesUsecase{
		typesRepository: or,
		timeout:       timeout,
	}
}

// Fetch returns questions with it's related repository
// Note: when using inmem repo, the usecase is a bit of useless
func (u *typesUsecase) Fetch(ctx context.Context, num int64) ([]*domain.Types, error) {
	return u.typesRepository.Fetch(ctx, num)
}

// GetByID returns single question based on given ID
func (u *typesUsecase) GetByID(ctx context.Context, id int) (*domain.Types, error) {
	return u.typesRepository.GetByID(ctx, id)
}

// Store stores question into repository
func (u *typesUsecase) Store(question domain.Types) (int64, error) {
	return u.typesRepository.Store(question)
}

// Update updates existing question into new question
func (u *typesUsecase) Update(question domain.Types) (domain.Types, error) {
	return u.typesRepository.Update(question)
}

// Delete removes existing question from repository based on given ID
func (u *typesUsecase) Delete(id int) error {
	return u.typesRepository.Delete(id)
}
