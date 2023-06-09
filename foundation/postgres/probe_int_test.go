package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/daniel-oliveiravas/class-booking-service/foundation/postgres"
)

func TestIntegrationPostgresProbe(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	ctx := context.Background()
	pool, err := postgres.Open(ctx, postgres.Config{
		Host:             "localhost",
		Port:             5432,
		DatabaseUser:     "class_booking",
		DatabasePassword: "class_booking",
		DatabaseName:     "class_booking_qa",
	})
	if err != nil {
		t.Fatalf("could not open database connection. err: %v", err)
	}

	postgresProbe := postgres.NewProbe(pool)

	err = postgresProbe.Check(ctx)
	if err != nil {
		t.Fatalf("it should ping database. err: %v", err)
	}

	t.Log("it should ping database.")
}
