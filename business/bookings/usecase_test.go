package bookings_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/class-booking-service/business/bookings"
	"github.com/daniel-oliveiravas/class-booking-service/business/bookings/mocks"
	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	classesmocks "github.com/daniel-oliveiravas/class-booking-service/business/classes/mocks"
	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	membersmocks "github.com/daniel-oliveiravas/class-booking-service/business/members/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_BookClass_MemberNotFound(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)

	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	now := time.Now()
	bookClass := bookings.BookClass{
		MemberID:  uuid.NewString(),
		ClassID:   uuid.NewString(),
		ClassDate: now,
	}

	expectedErr := pgx.ErrNoRows
	membersRepo.On("GetByID", mock.Anything, bookClass.MemberID).Return(members.Member{}, expectedErr).Once()
	membersRepo.On("IsNotFoundErr", expectedErr).Return(true).Once()

	_, err := usecase.BookClass(ctx, bookClass)
	require.Error(t, err)
	require.True(t, errors.Is(err, bookings.ErrMemberNotFound))
}

func TestUsecase_BookClass_ClassNotFound(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)

	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	now := time.Now()
	bookClass := bookings.BookClass{
		MemberID:  uuid.NewString(),
		ClassID:   uuid.NewString(),
		ClassDate: now,
	}

	membersRepo.On("GetByID", mock.Anything, bookClass.MemberID).Return(members.Member{ID: bookClass.MemberID}, nil).Once()

	expectedErr := pgx.ErrNoRows
	classesRepo.On("GetByID", mock.Anything, bookClass.ClassID).Return(classes.Class{}, expectedErr).Once()
	classesRepo.On("IsNotFoundErr", expectedErr).Return(true).Once()

	_, err := usecase.BookClass(ctx, bookClass)
	require.Error(t, err)
	require.True(t, errors.Is(err, bookings.ErrClassNotFound))
}

func TestUsecase_BookClass_DateBeforeStartClass(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)

	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	now := time.Now()
	bookClass := bookings.BookClass{
		MemberID:  uuid.NewString(),
		ClassID:   uuid.NewString(),
		ClassDate: now,
	}

	membersRepo.On("GetByID", mock.Anything, bookClass.MemberID).Return(members.Member{ID: bookClass.MemberID}, nil).Once()
	classesRepo.On("GetByID", mock.Anything, bookClass.ClassID).Return(classes.Class{StartDate: now.Add(time.Hour)}, nil).Once()

	_, err := usecase.BookClass(ctx, bookClass)
	require.Error(t, err)
	require.True(t, errors.Is(err, bookings.ErrInvalidClassDate))
}

func TestUsecase_BookClass_DateAfterEndClass(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)

	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	now := time.Now()
	bookClass := bookings.BookClass{
		MemberID:  uuid.NewString(),
		ClassID:   uuid.NewString(),
		ClassDate: now,
	}

	membersRepo.On("GetByID", mock.Anything, bookClass.MemberID).Return(members.Member{ID: bookClass.MemberID}, nil).Once()
	classesRepo.On("GetByID", mock.Anything, bookClass.ClassID).Return(classes.Class{StartDate: now.Add(time.Hour * -2), EndDate: now.Add(time.Hour * -1)}, nil).Once()

	_, err := usecase.BookClass(ctx, bookClass)
	require.Error(t, err)
	require.True(t, errors.Is(err, bookings.ErrInvalidClassDate))
}

func TestUsecase_BookClass_ValidDate(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)

	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	now := time.Now()
	bookClass := bookings.BookClass{
		MemberID:  uuid.NewString(),
		ClassID:   uuid.NewString(),
		ClassDate: now,
	}

	membersRepo.On("GetByID", mock.Anything, bookClass.MemberID).Return(members.Member{ID: bookClass.MemberID}, nil).Once()
	classesRepo.On("GetByID", mock.Anything, bookClass.ClassID).Return(classes.Class{StartDate: now.Add(time.Hour * -2), EndDate: now.Add(time.Hour * 1)}, nil).Once()
	repo.On("BookClass", mock.Anything, mock.Anything).Return(NewBooking(), nil).Once()

	_, err := usecase.BookClass(ctx, bookClass)
	require.NoError(t, err)
}

func NewBooking() bookings.Booking {
	return bookings.Booking{
		ID:        uuid.NewString(),
		MemberID:  uuid.NewString(),
		ClassID:   uuid.NewString(),
		ClassDate: time.Time{},
		BookedAt:  time.Time{},
		UpdatedAt: time.Time{},
	}
}

func TestUsecase_GetByID(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)
	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	bookingID := uuid.NewString()
	expectedClass := NewBooking()

	repo.On("GetByID", mock.Anything, bookingID).Return(expectedClass, nil).Once()
	bookingFound, err := usecase.GetByID(ctx, bookingID)
	require.NoError(t, err)
	assert.Equal(t, expectedClass, bookingFound)
}

func TestUsecase_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)
	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	bookingID := uuid.NewString()

	expectedErr := pgx.ErrNoRows
	repo.On("GetByID", mock.Anything, bookingID).Return(bookings.Booking{}, expectedErr).Once()
	repo.On("IsNotFoundErr", expectedErr).Return(true).Once()
	bookingFound, err := usecase.GetByID(ctx, bookingID)
	require.Error(t, err)
	require.Empty(t, bookingFound.ID)
	require.True(t, errors.Is(err, bookings.ErrNotFound))
}

func TestUsecase_DeleteBooking(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)
	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	bookingID := uuid.NewString()
	repo.On("DeleteBooking", mock.Anything, bookingID).Return(nil).Once()
	err := usecase.DeleteBooking(ctx, bookingID)
	require.NoError(t, err)
}

func TestUsecase_ListBookings(t *testing.T) {
	ctx := context.Background()
	membersRepo := membersmocks.NewRepository(t)
	membersUsecase := members.NewUsecase(membersRepo)
	classesRepo := classesmocks.NewRepository(t)
	classesUsecase := classes.NewUsecase(classesRepo)
	repo := mocks.NewRepository(t)
	usecase := bookings.NewUsecase(repo, membersUsecase, classesUsecase)

	expectedBookings := make([]bookings.Booking, 0)
	expectedBookings = append(expectedBookings, NewBooking(), NewBooking())
	repo.On("ListBookings", mock.Anything, 100, 0).Return(expectedBookings, nil).Once()

	pageInfo := bookings.PageInfo{
		Limit: 200,
		Page:  0,
	}
	all, err := usecase.ListBookings(ctx, pageInfo)
	require.NoError(t, err)
	require.NotEmpty(t, all)
}
