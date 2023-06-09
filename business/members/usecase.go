package members

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidData = errors.New("invalid data")
	ErrNotFound    = errors.New("not found")
)

type Usecase struct {
	repository Repository
}

func NewUsecase(repository Repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

//go:generate mockery --name=Repository --filename=members_repository.go
type Repository interface {
	AddMember(ctx context.Context, member Member) (Member, error)
	GetByID(ctx context.Context, memberID string) (Member, error)
	IsNotFoundErr(err error) bool
	UpdateMember(ctx context.Context, memberID string, updateMember UpdateMember) (Member, error)
	DeleteMember(ctx context.Context, memberID string) error
	ListMembers(ctx context.Context, limit int, offset int) ([]Member, error)
}

func (u *Usecase) AddMember(ctx context.Context, newMember NewMember) (Member, error) {
	member := Member{
		ID:   uuid.NewString(),
		Name: newMember.Name,
	}

	if err := u.validateMember(member); err != nil {
		return Member{}, err
	}

	addedMember, err := u.repository.AddMember(ctx, member)
	if err != nil {
		return Member{}, fmt.Errorf("failed to add member to repository: %w", err)
	}

	return addedMember, nil
}

func (u *Usecase) GetByID(ctx context.Context, memberID string) (Member, error) {
	member, err := u.repository.GetByID(ctx, memberID)
	if err != nil {
		if u.repository.IsNotFoundErr(err) {
			return Member{}, ErrNotFound
		}

		return Member{}, err
	}

	return member, nil
}

func (u *Usecase) UpdateMember(ctx context.Context, memberID string, updateMember UpdateMember) (Member, error) {
	updatedMember, err := u.repository.UpdateMember(ctx, memberID, updateMember)
	if err != nil {
		if u.repository.IsNotFoundErr(err) {
			return Member{}, ErrNotFound
		}
		return Member{}, fmt.Errorf("failed to update member in repository: %w", err)
	}

	return updatedMember, nil
}

func (u *Usecase) DeleteMember(ctx context.Context, memberID string) error {
	return u.repository.DeleteMember(ctx, memberID)
}

func (u *Usecase) ListMembers(ctx context.Context, pageInfo PageInfo) ([]Member, error) {
	if pageInfo.Limit > 100 || pageInfo.Limit == 0 {
		pageInfo.Limit = 100
	}

	offset := pageInfo.Limit * pageInfo.Page
	return u.repository.ListMembers(ctx, pageInfo.Limit, offset)
}

func (u *Usecase) validateMember(member Member) error {
	if member.Name == "" {
		return fmt.Errorf("missing member name. %w", ErrInvalidData)
	}

	return nil
}
