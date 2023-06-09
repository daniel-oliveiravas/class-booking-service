package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/class-booking-service/business/bookings"
	pgrepo "github.com/daniel-oliveiravas/class-booking-service/business/bookings/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	pgrepoclass "github.com/daniel-oliveiravas/class-booking-service/business/classes/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	pgrepomember "github.com/daniel-oliveiravas/class-booking-service/business/members/integration/postgres"
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

func TestRepository_AddBooking(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop().Sugar()

	repo := pgrepo.NewBookingsRepository(logger, db)
	memberRepo := pgrepomember.NewMembersRepository(logger, db)
	classRepo := pgrepoclass.NewClassesRepository(logger, db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	memberAdded, err := memberRepo.AddMember(ctx, member)
	require.NoError(t, err)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	classAdded, err := classRepo.Add(ctx, class)
	require.NoError(t, err)

	classDate := time.Now().UTC()
	bookClass := bookings.Booking{
		ID:        uuid.NewString(),
		MemberID:  memberAdded.ID,
		ClassID:   classAdded.ID,
		ClassDate: classDate,
	}
	booking, err := repo.BookClass(ctx, bookClass)
	require.NoError(t, err)

	assert.Equal(t, bookClass.ID, booking.ID)
	assert.Equal(t, bookClass.MemberID, booking.MemberID)
	assert.Equal(t, bookClass.ClassID, booking.ClassID)
	assert.NotEmpty(t, booking.BookedAt)
	assert.NotEmpty(t, booking.UpdatedAt)
}

func TestRepository_GetByID(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop().Sugar()

	repo := pgrepo.NewBookingsRepository(logger, db)
	memberRepo := pgrepomember.NewMembersRepository(logger, db)
	classRepo := pgrepoclass.NewClassesRepository(logger, db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	memberAdded, err := memberRepo.AddMember(ctx, member)
	require.NoError(t, err)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	classAdded, err := classRepo.Add(ctx, class)
	require.NoError(t, err)

	classDate := time.Now().UTC()
	bookClass := bookings.Booking{
		ID:        uuid.NewString(),
		MemberID:  memberAdded.ID,
		ClassID:   classAdded.ID,
		ClassDate: classDate,
	}
	booking, err := repo.BookClass(ctx, bookClass)
	require.NoError(t, err)

	bookingFound, err := repo.GetByID(ctx, booking.ID)
	require.NoError(t, err)

	assert.Equal(t, booking.ID, bookingFound.ID)
	assert.Equal(t, booking.MemberID, bookingFound.MemberID)
	assert.Equal(t, booking.ClassID, bookingFound.ClassID)
	assert.Equal(t, booking.ClassDate, bookingFound.ClassDate)
	assert.NotEmpty(t, bookingFound.BookedAt)
	assert.NotEmpty(t, bookingFound.UpdatedAt)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop()

	repo := pgrepo.NewBookingsRepository(logger.Sugar(), db)

	_, err := repo.GetByID(ctx, uuid.NewString())
	assert.True(t, repo.IsNotFoundErr(err))
}

func TestRepository_DeleteMember(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop().Sugar()

	repo := pgrepo.NewBookingsRepository(logger, db)
	memberRepo := pgrepomember.NewMembersRepository(logger, db)
	classRepo := pgrepoclass.NewClassesRepository(logger, db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	memberAdded, err := memberRepo.AddMember(ctx, member)
	require.NoError(t, err)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	classAdded, err := classRepo.Add(ctx, class)
	require.NoError(t, err)

	classDate := time.Now().UTC()
	bookClass := bookings.Booking{
		ID:        uuid.NewString(),
		MemberID:  memberAdded.ID,
		ClassID:   classAdded.ID,
		ClassDate: classDate,
	}
	booking, err := repo.BookClass(ctx, bookClass)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, booking.ID)
	require.False(t, repo.IsNotFoundErr(err))

	err = repo.DeleteBooking(ctx, booking.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, booking.ID)
	require.True(t, repo.IsNotFoundErr(err))
}

func TestRepository_ListBookings(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	db := setupIntegration(t)
	ctx := context.Background()
	logger := zap.NewNop().Sugar()

	repo := pgrepo.NewBookingsRepository(logger, db)
	memberRepo := pgrepomember.NewMembersRepository(logger, db)
	classRepo := pgrepoclass.NewClassesRepository(logger, db)

	member := members.Member{
		ID:   uuid.NewString(),
		Name: uuid.NewString(),
	}
	memberAdded, err := memberRepo.AddMember(ctx, member)
	require.NoError(t, err)

	now := time.Now().UTC()
	class := classes.Class{
		ID:        uuid.NewString(),
		Name:      uuid.NewString(),
		StartDate: now,
		EndDate:   now,
		Capacity:  20,
	}
	classAdded, err := classRepo.Add(ctx, class)
	require.NoError(t, err)

	classDate := time.Now().UTC()
	bookClass := bookings.Booking{
		ID:        uuid.NewString(),
		MemberID:  memberAdded.ID,
		ClassID:   classAdded.ID,
		ClassDate: classDate,
	}
	_, err = repo.BookClass(ctx, bookClass)
	require.NoError(t, err)

	allMembers, err := repo.ListBookings(ctx, 100, 0)
	require.NoError(t, err)
	require.NotEmpty(t, allMembers)
}
