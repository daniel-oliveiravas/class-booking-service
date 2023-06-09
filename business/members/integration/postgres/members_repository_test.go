package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	pgrepo "github.com/daniel-oliveiravas/class-booking-service/business/members/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/foundation/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupIntegration(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	schema := t.Name()
	pgCfg := postgres.Config{
		Host:             "localhost",
		Port:             5432,
		DatabaseUser:     "class_booking",
		DatabasePassword: "class_booking",
		DatabaseName:     "class_booking_qa",
		SSLMode:          "none",
		SearchPath:       schema,
	}
	db, err := postgres.Open(ctx, pgCfg)
	require.NoError(t, err)

	err = postgres.DropAndCreateSchema(ctx, db, schema)
	require.NoError(t, err)

	err = postgres.Migrate("file://../../../../scripts/db/migrations/", pgCfg)
	require.NoError(t, err)

	return db
}

func TestRepository_AddMember(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewMembersRepository(logger.Sugar(), db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	memberAdded, err := repo.AddMember(ctx, member)
	require.NoError(t, err)

	assert.Equal(t, member.ID, memberAdded.ID)
	assert.Equal(t, member.Name, memberAdded.Name)
	assert.NotEmpty(t, memberAdded.CreatedAt)
	assert.NotEmpty(t, memberAdded.UpdatedAt)
}

func TestRepository_GetByID(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewMembersRepository(logger.Sugar(), db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	memberAdded, err := repo.AddMember(ctx, member)
	require.NoError(t, err)

	memberFound, err := repo.GetByID(ctx, memberAdded.ID)
	require.NoError(t, err)

	assert.Equal(t, memberAdded.ID, memberFound.ID)
	assert.Equal(t, memberAdded.Name, memberFound.Name)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewMembersRepository(logger.Sugar(), db)

	_, err := repo.GetByID(ctx, uuid.NewString())
	assert.True(t, repo.IsNotFoundErr(err))
}

func TestRepository_UpdateMember(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewMembersRepository(logger.Sugar(), db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	_, err := repo.AddMember(ctx, member)
	require.NoError(t, err)

	newName := uuid.NewString()
	updatedMember := members.UpdateMember{
		Name: &newName,
	}

	memberUpdated, err := repo.UpdateMember(ctx, member.ID, updatedMember)
	require.NoError(t, err)

	assert.Equal(t, member.ID, memberUpdated.ID)
	assert.Equal(t, newName, memberUpdated.Name)
}

func TestRepository_DeleteMember(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewMembersRepository(logger.Sugar(), db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	_, err := repo.AddMember(ctx, member)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, member.ID)
	require.False(t, repo.IsNotFoundErr(err))

	err = repo.DeleteMember(ctx, member.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, member.ID)
	require.True(t, repo.IsNotFoundErr(err))
}

func TestRepository_ListMembers(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewMembersRepository(logger.Sugar(), db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	_, err := repo.AddMember(ctx, member)
	require.NoError(t, err)

	allMembers, err := repo.ListMembers(ctx, 100, 0)
	require.NoError(t, err)
	require.NotEmpty(t, allMembers)
}
