package classes

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

//go:generate mockery --name=Repository --filename=classes_repository.go
type Repository interface {
	Add(ctx context.Context, class Class) (Class, error)
	GetByID(ctx context.Context, classID string) (Class, error)
	IsNotFoundErr(err error) bool
	Update(ctx context.Context, classID string, updateClass UpdateClass) (Class, error)
	Delete(ctx context.Context, classID string) error
	List(ctx context.Context, limit int, offset int) ([]Class, error)
}

func (u *Usecase) AddClass(ctx context.Context, newClass NewClass) (Class, error) {
	class := Class{
		ID:        uuid.NewString(),
		Name:      newClass.Name,
		StartDate: newClass.StartDate,
		EndDate:   newClass.EndDate,
		Capacity:  newClass.Capacity,
	}

	if err := u.validClass(class); err != nil {
		return Class{}, err
	}

	classAdded, err := u.repository.Add(ctx, class)
	if err != nil {
		return Class{}, fmt.Errorf("failed to add class to repository: %w", err)
	}

	return classAdded, nil
}

func (u *Usecase) GetByID(ctx context.Context, classID string) (Class, error) {
	class, err := u.repository.GetByID(ctx, classID)
	if err != nil {
		if u.repository.IsNotFoundErr(err) {
			return Class{}, ErrNotFound
		}

		return Class{}, err
	}

	return class, nil
}

func (u *Usecase) UpdateClass(ctx context.Context, classID string, updateClass UpdateClass) (Class, error) {
	classUpdated, err := u.repository.Update(ctx, classID, updateClass)
	if err != nil {
		if u.repository.IsNotFoundErr(err) {
			return Class{}, ErrNotFound
		}
		return Class{}, fmt.Errorf("failed to update class in repository: %w", err)
	}

	return classUpdated, nil
}

func (u *Usecase) DeleteClass(ctx context.Context, classID string) error {
	return u.repository.Delete(ctx, classID)
}

func (u *Usecase) ListClasses(ctx context.Context, pageInfo PageInfo) ([]Class, error) {
	if pageInfo.Limit > 100 || pageInfo.Limit == 0 {
		pageInfo.Limit = 100
	}

	offset := pageInfo.Limit * pageInfo.Page
	return u.repository.List(ctx, pageInfo.Limit, offset)
}

func (u *Usecase) validClass(class Class) error {
	if class.Name == "" {
		return fmt.Errorf("missing class 'name': %w", ErrInvalidData)
	}

	if class.Capacity == 0 {
		return fmt.Errorf("missing class 'capacity': %w", ErrInvalidData)
	}

	if class.StartDate.After(class.EndDate) {
		return fmt.Errorf("start date cannot be later than end date: %w", ErrInvalidData)
	}

	return nil
}
