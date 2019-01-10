package usecase

import (
	"context"
	"time"

	"github.com/DynEd/etest/domain"
)

type groupsUsecase struct {
	groupsRepository domain.GroupsRepository
	timeout            time.Duration
}

// NewGroupsUsecase will create new an questionUsecase object representation of domain.GroupsUsecase interface
func NewGroupsUsecase(or domain.GroupsRepository, timeout time.Duration) domain.GroupsUsecase {
	return &groupsUsecase{
		groupsRepository: or,
		timeout:       timeout,
	}
}

// Fetch returns questions with it's related repository
// Note: when using inmem repo, the usecase is a bit of useless
func (u *groupsUsecase) Fetch(ctx context.Context, num int64) ([]*domain.Groups, error) {
	return u.groupsRepository.Fetch(ctx, num)
}

// GetByID returns single question based on given ID
func (u *groupsUsecase) GetByID(ctx context.Context, id int) (*domain.Groups, error) {
	return u.groupsRepository.GetByID(ctx, id)
}

// Store stores question into repository
func (u *groupsUsecase) Store(question domain.Groups) (int64, error) {
	return u.groupsRepository.Store(question)
}

// Update updates existing question into new question
func (u *groupsUsecase) Update(question domain.Groups) (domain.Groups, error) {
	return u.groupsRepository.Update(question)
}

// Delete removes existing question from repository based on given ID
func (u *groupsUsecase) Delete(id int) error {
	return u.groupsRepository.Delete(id)
}
