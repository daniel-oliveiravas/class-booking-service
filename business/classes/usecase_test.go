package classes_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	"github.com/daniel-oliveiravas/class-booking-service/business/classes/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_AddClass(t *testing.T) {
	classesRepo := mocks.NewRepository(t)
	usecase := classes.NewUsecase(classesRepo)
	tests := []struct {
		name     string
		newClass classes.NewClass
		want     classes.Class
		wantErr  error
	}{
		{
			name: "without_name",
			newClass: classes.NewClass{
				Name: "",
			},
			want:    classes.Class{},
			wantErr: classes.ErrInvalidData,
		},
		{
			name: "without_capacity",
			newClass: classes.NewClass{
				Name:     uuid.NewString(),
				Capacity: 0,
			},
			want:    classes.Class{},
			wantErr: classes.ErrInvalidData,
		},
		{
			name: "start_date_after_end_date",
			newClass: classes.NewClass{
				Name:      uuid.NewString(),
				Capacity:  30,
				StartDate: time.Now().Add(time.Hour * 24),
				EndDate:   time.Now(),
			},
			want:    classes.Class{},
			wantErr: classes.ErrInvalidData,
		},
		{
			name: "valid_class",
			newClass: classes.NewClass{
				Name:      uuid.NewString(),
				Capacity:  30,
				StartDate: time.Now(),
				EndDate:   time.Now().Add(time.Hour * 24 * 10),
			},
			want:    NewClass(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.wantErr == nil {
				classesRepo.On("Add", mock.Anything, mock.Anything).Return(tt.want, nil).Once()
			}

			class, err := usecase.AddClass(ctx, tt.newClass)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.wantErr))
				return
			}

			assert.Equal(t, tt.want, class)
		})
	}
}

func NewClass() classes.Class {
	return classes.Class{
		ID:        uuid.NewString(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Name:      uuid.NewString(),
		StartDate: time.Time{},
		EndDate:   time.Time{},
		Capacity:  30,
	}
}

func TestUsecase_GetByID(t *testing.T) {
	ctx := context.Background()
	classesRepo := mocks.NewRepository(t)
	usecase := classes.NewUsecase(classesRepo)

	classID := uuid.NewString()
	expectedClass := NewClass()

	classesRepo.On("GetByID", mock.Anything, classID).Return(expectedClass, nil).Once()
	classFound, err := usecase.GetByID(ctx, classID)
	require.NoError(t, err)
	assert.Equal(t, expectedClass, classFound)
}

func TestUsecase_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	classesRepo := mocks.NewRepository(t)
	usecase := classes.NewUsecase(classesRepo)

	classID := uuid.NewString()

	expectedErr := pgx.ErrNoRows
	classesRepo.On("GetByID", mock.Anything, classID).Return(classes.Class{}, expectedErr).Once()
	classesRepo.On("IsNotFoundErr", expectedErr).Return(true).Once()
	classFound, err := usecase.GetByID(ctx, classID)
	require.Error(t, err)
	require.Empty(t, classFound.ID)
	require.True(t, errors.Is(err, classes.ErrNotFound))
}

func TestUsecase_UpdateClass_NotFound(t *testing.T) {
	ctx := context.Background()
	classesRepo := mocks.NewRepository(t)
	usecase := classes.NewUsecase(classesRepo)

	classID := uuid.NewString()

	UpdateClass := classes.UpdateClass{}
	expectedErr := pgx.ErrNoRows
	classesRepo.On("Update", mock.Anything, classID, UpdateClass).Return(classes.Class{}, expectedErr).Once()
	classesRepo.On("IsNotFoundErr", expectedErr).Return(true).Once()

	classFound, err := usecase.UpdateClass(ctx, classID, UpdateClass)
	require.Error(t, err)
	require.Empty(t, classFound.ID)
	require.True(t, errors.Is(err, classes.ErrNotFound))
}

func TestUsecase_UpdateClass(t *testing.T) {
	ctx := context.Background()
	classesRepo := mocks.NewRepository(t)
	usecase := classes.NewUsecase(classesRepo)

	classID := uuid.NewString()

	expectedClass := NewClass()
	newName := uuid.NewString()
	UpdateClass := classes.UpdateClass{
		Name: &newName,
	}

	expectedClass.Name = newName
	classesRepo.On("Update", mock.Anything, classID, UpdateClass).Return(expectedClass, nil).Once()

	classUpdated, err := usecase.UpdateClass(ctx, classID, UpdateClass)
	require.NoError(t, err)
	require.Equal(t, newName, classUpdated.Name)
}

func TestUsecase_DeleteClass(t *testing.T) {
	ctx := context.Background()
	classesRepo := mocks.NewRepository(t)
	usecase := classes.NewUsecase(classesRepo)

	classID := uuid.NewString()
	classesRepo.On("Delete", mock.Anything, classID).Return(nil).Once()
	err := usecase.DeleteClass(ctx, classID)
	require.NoError(t, err)
}

func TestUsecase_ListClasses(t *testing.T) {
	ctx := context.Background()
	classesRepo := mocks.NewRepository(t)
	usecase := classes.NewUsecase(classesRepo)

	expectedClasses := make([]classes.Class, 0)
	expectedClasses = append(expectedClasses, NewClass(), NewClass())
	classesRepo.On("List", mock.Anything, 100, 0).Return(expectedClasses, nil).Once()

	pageInfo := classes.PageInfo{
		Limit: 200,
		Page:  0,
	}
	all, err := usecase.ListClasses(ctx, pageInfo)
	require.NoError(t, err)
	require.NotEmpty(t, all)
}
