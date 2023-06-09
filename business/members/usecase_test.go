package members_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	"github.com/daniel-oliveiravas/class-booking-service/business/members/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_AddMember(t *testing.T) {
	membersRepo := mocks.NewRepository(t)
	usecase := members.NewUsecase(membersRepo)
	tests := []struct {
		name      string
		newMember members.NewMember
		want      members.Member
		wantErr   error
	}{
		{
			name: "without_name",
			newMember: members.NewMember{
				Name: "",
			},
			want:    members.Member{},
			wantErr: members.ErrInvalidData,
		},
		{
			name: "valid_member",
			newMember: members.NewMember{
				Name: uuid.NewString(),
			},
			want:    NewMember(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.wantErr == nil {
				membersRepo.On("AddMember", mock.Anything, mock.Anything).Return(tt.want, nil).Once()
			}

			member, err := usecase.AddMember(ctx, tt.newMember)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.wantErr))
				return
			}

			assert.Equal(t, tt.want, member)
		})
	}
}

func NewMember() members.Member {
	return members.Member{
		ID:        uuid.NewString(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Name:      uuid.NewString(),
	}
}

func TestUsecase_GetByID(t *testing.T) {
	ctx := context.Background()
	membersRepo := mocks.NewRepository(t)
	usecase := members.NewUsecase(membersRepo)

	memberID := uuid.NewString()
	expectedMember := NewMember()

	membersRepo.On("GetByID", mock.Anything, memberID).Return(expectedMember, nil).Once()
	memberFound, err := usecase.GetByID(ctx, memberID)
	require.NoError(t, err)
	assert.Equal(t, expectedMember, memberFound)
}

func TestUsecase_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	membersRepo := mocks.NewRepository(t)
	usecase := members.NewUsecase(membersRepo)

	memberID := uuid.NewString()

	expectedErr := pgx.ErrNoRows
	membersRepo.On("GetByID", mock.Anything, memberID).Return(members.Member{}, expectedErr).Once()
	membersRepo.On("IsNotFoundErr", expectedErr).Return(true).Once()
	memberFound, err := usecase.GetByID(ctx, memberID)
	require.Error(t, err)
	require.Empty(t, memberFound.ID)
	require.True(t, errors.Is(err, members.ErrNotFound))
}

func TestUsecase_UpdateMember_NotFound(t *testing.T) {
	ctx := context.Background()
	membersRepo := mocks.NewRepository(t)
	usecase := members.NewUsecase(membersRepo)

	memberID := uuid.NewString()

	updateMember := members.UpdateMember{}
	expectedErr := pgx.ErrNoRows
	membersRepo.On("UpdateMember", mock.Anything, memberID, updateMember).Return(members.Member{}, expectedErr).Once()
	membersRepo.On("IsNotFoundErr", expectedErr).Return(true).Once()

	memberFound, err := usecase.UpdateMember(ctx, memberID, updateMember)
	require.Error(t, err)
	require.Empty(t, memberFound.ID)
	require.True(t, errors.Is(err, members.ErrNotFound))
}

func TestUsecase_UpdateMember(t *testing.T) {
	ctx := context.Background()
	membersRepo := mocks.NewRepository(t)
	usecase := members.NewUsecase(membersRepo)

	memberID := uuid.NewString()

	expectedMember := NewMember()
	newName := uuid.NewString()
	updateMember := members.UpdateMember{
		Name: &newName,
	}

	expectedMember.Name = newName
	membersRepo.On("UpdateMember", mock.Anything, memberID, updateMember).Return(expectedMember, nil).Once()

	memberUpdated, err := usecase.UpdateMember(ctx, memberID, updateMember)
	require.NoError(t, err)
	require.Equal(t, newName, memberUpdated.Name)
}

func TestUsecase_DeleteMember(t *testing.T) {
	ctx := context.Background()
	membersRepo := mocks.NewRepository(t)
	usecase := members.NewUsecase(membersRepo)

	memberID := uuid.NewString()
	membersRepo.On("DeleteMember", mock.Anything, memberID).Return(nil).Once()
	err := usecase.DeleteMember(ctx, memberID)
	require.NoError(t, err)
}

func TestUsecase_ListMembers(t *testing.T) {
	ctx := context.Background()
	membersRepo := mocks.NewRepository(t)
	usecase := members.NewUsecase(membersRepo)

	expectedMembers := make([]members.Member, 0)
	expectedMembers = append(expectedMembers, NewMember(), NewMember())
	membersRepo.On("ListMembers", mock.Anything, 100, 0).Return(expectedMembers, nil).Once()

	pageInfo := members.PageInfo{
		Limit: 200,
		Page:  0,
	}
	all, err := usecase.ListMembers(ctx, pageInfo)
	require.NoError(t, err)
	require.NotEmpty(t, all)
}
