package bookings

import (
	"context"
	"errors"
	"fmt"

	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	"github.com/google/uuid"
)

var (
	ErrInvalidClassDate = errors.New("invalid class date")
	ErrNotFound         = errors.New("booking not found")
	ErrMemberNotFound   = errors.New("member not found")
	ErrClassNotFound    = errors.New("class not found")
)

type Usecase struct {
	repository     Repository
	membersUsecase *members.Usecase
	classesUsecase *classes.Usecase
}

func NewUsecase(repository Repository, membersUsecase *members.Usecase, classesUsecase *classes.Usecase) *Usecase {
	return &Usecase{repository: repository,
		membersUsecase: membersUsecase,
		classesUsecase: classesUsecase,
	}
}

//go:generate mockery --name=Repository --filename=booking_repository.go
type Repository interface {
	BookClass(ctx context.Context, booking Booking) (Booking, error)
	GetByID(ctx context.Context, bookingID string) (Booking, error)
	IsNotFoundErr(err error) bool
	DeleteBooking(ctx context.Context, bookingID string) error
	ListBookings(ctx context.Context, limit int, offset int) ([]Booking, error)
}

func (u *Usecase) BookClass(ctx context.Context, bookClass BookClass) (Booking, error) {
	booking := Booking{
		ID:        uuid.NewString(),
		MemberID:  bookClass.MemberID,
		ClassID:   bookClass.ClassID,
		ClassDate: bookClass.ClassDate,
	}

	if err := u.validateBooking(ctx, booking); err != nil {
		return Booking{}, err
	}

	classAdded, err := u.repository.BookClass(ctx, booking)
	if err != nil {
		return Booking{}, fmt.Errorf("failed to add booking to repository: %w", err)
	}

	return classAdded, nil
}

func (u *Usecase) GetByID(ctx context.Context, classID string) (Booking, error) {
	class, err := u.repository.GetByID(ctx, classID)
	if err != nil {
		if u.repository.IsNotFoundErr(err) {
			return Booking{}, ErrNotFound
		}

		return Booking{}, err
	}

	return class, nil
}

func (u *Usecase) DeleteBooking(ctx context.Context, bookingID string) error {
	return u.repository.DeleteBooking(ctx, bookingID)
}

func (u *Usecase) ListBookings(ctx context.Context, pageInfo PageInfo) ([]Booking, error) {
	if pageInfo.Limit > 100 || pageInfo.Limit == 0 {
		pageInfo.Limit = 100
	}

	offset := pageInfo.Limit * pageInfo.Page
	return u.repository.ListBookings(ctx, pageInfo.Limit, offset)
}

func (u *Usecase) validateBooking(ctx context.Context, booking Booking) error {
	_, err := u.membersUsecase.GetByID(ctx, booking.MemberID)
	if err != nil {
		if errors.Is(err, members.ErrNotFound) {
			return ErrMemberNotFound
		}

		return err
	}

	class, err := u.classesUsecase.GetByID(ctx, booking.ClassID)
	if err != nil {
		if errors.Is(err, classes.ErrNotFound) {
			return ErrClassNotFound
		}
	}

	if booking.ClassDate.Before(class.StartDate) || booking.ClassDate.After(class.EndDate) {
		return ErrInvalidClassDate
	}

	return nil
}
