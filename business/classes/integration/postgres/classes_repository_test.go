package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	pgrepo "github.com/daniel-oliveiravas/class-booking-service/business/classes/integration/postgres"
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

func TestRepository_AddClass(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewClassesRepository(logger.Sugar(), db)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	addedClass, err := repo.Add(ctx, class)
	require.NoError(t, err)

	assert.Equal(t, class.ID, addedClass.ID)
	assert.NotEmpty(t, addedClass.CreatedAt)
	assert.NotEmpty(t, addedClass.UpdatedAt)
	assert.Equal(t, class.Name, addedClass.Name)
	assert.Equal(t, class.StartDate, addedClass.StartDate)
	assert.Equal(t, class.EndDate, addedClass.EndDate)
	assert.Equal(t, class.Capacity, addedClass.Capacity)
}

func TestRepository_GetByID(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewClassesRepository(logger.Sugar(), db)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	_, err := repo.Add(ctx, class)
	require.NoError(t, err)

	classFound, err := repo.GetByID(ctx, class.ID)
	require.NoError(t, err)

	assert.Equal(t, class.ID, classFound.ID)
	assert.NotEmpty(t, classFound.CreatedAt)
	assert.NotEmpty(t, classFound.UpdatedAt)
	assert.Equal(t, class.Name, classFound.Name)
	assert.Equal(t, class.StartDate, classFound.StartDate)
	assert.Equal(t, class.EndDate, classFound.EndDate)
	assert.Equal(t, class.Capacity, classFound.Capacity)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewClassesRepository(logger.Sugar(), db)

	_, err := repo.GetByID(ctx, uuid.NewString())
	assert.True(t, repo.IsNotFoundErr(err))
}

func TestRepository_UpdateClass(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewClassesRepository(logger.Sugar(), db)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	_, err := repo.Add(ctx, class)
	require.NoError(t, err)

	newCapacity := 45
	newName := uuid.NewString()
	updateClass := classes.UpdateClass{
		Name:       &newName,
		Capability: &newCapacity,
	}

	updatedClass, err := repo.Update(ctx, class.ID, updateClass)
	require.NoError(t, err)

	assert.Equal(t, class.ID, updatedClass.ID)
	assert.Equal(t, newName, updatedClass.Name)
	assert.Equal(t, newCapacity, updatedClass.Capacity)
}

func TestRepository_DeleteClass(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewClassesRepository(logger.Sugar(), db)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	_, err := repo.Add(ctx, class)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, class.ID)
	require.False(t, repo.IsNotFoundErr(err))

	err = repo.Delete(ctx, class.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, class.ID)
	require.True(t, repo.IsNotFoundErr(err))
}

func TestRepository_ListClass(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewClassesRepository(logger.Sugar(), db)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	_, err := repo.Add(ctx, class)
	require.NoError(t, err)

	allClasses, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.NotEmpty(t, allClasses)
}
