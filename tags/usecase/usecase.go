package usecase

import (
	"context"
	"time"

	"github.com/DynEd/etest/domain"
)

type tagsUsecase struct {
	tagsRepository domain.TagsRepository
	timeout            time.Duration
}

// NewTagsUsecase will create new an questionUsecase object representation of domain.TagsUsecase interface
func NewTagsUsecase(or domain.TagsRepository, timeout time.Duration) domain.TagsUsecase {
	return &tagsUsecase{
		tagsRepository: or,
		timeout:       timeout,
	}
}

// Fetch returns questions with it's related repository
// Note: when using inmem repo, the usecase is a bit of useless
func (u *tagsUsecase) Fetch(ctx context.Context, num int64) ([]*domain.Tags, error) {
	return u.tagsRepository.Fetch(ctx, num)
}

// GetByID returns single question based on given ID
func (u *tagsUsecase) GetByID(ctx context.Context, id int) (*domain.Tags, error) {
	return u.tagsRepository.GetByID(ctx, id)
}

// Store stores question into repository
func (u *tagsUsecase) Store(question domain.Tags) (int64, error) {
	return u.tagsRepository.Store(question)
}

// Update updates existing question into new question
func (u *tagsUsecase) Update(question domain.Tags) (domain.Tags, error) {
	return u.tagsRepository.Update(question)
}

// Delete removes existing question from repository based on given ID
func (u *tagsUsecase) Delete(id int) error {
	return u.tagsRepository.Delete(id)
}
