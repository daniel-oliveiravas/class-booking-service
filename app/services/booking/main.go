package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/daniel-oliveiravas/class-booking-service/app/services/booking/handlers"
	"github.com/daniel-oliveiravas/class-booking-service/business/bookings"
	pgbookings "github.com/daniel-oliveiravas/class-booking-service/business/bookings/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	pgclasses "github.com/daniel-oliveiravas/class-booking-service/business/classes/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	pgmembers "github.com/daniel-oliveiravas/class-booking-service/business/members/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/foundation/logging"
	"github.com/daniel-oliveiravas/class-booking-service/foundation/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func main() {
	logger, err := logging.New("class-booking-api")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer logger.Sync()

	if err := run(logger); err != nil {
		logger.Errorw("failed to start service", "error", err.Error())
	}
}

func run(logger *zap.SugaredLogger) error {
	ctx := context.Background()

	// -------------------------------------------------------------------------
	// Configuration
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// -------------------------------------------------------------------------
	// Starting Database Connection

	pgCfg := postgres.Config{
		Host:             cfg.PostgresHostname,
		Port:             cfg.PostgresPort,
		DatabaseUser:     cfg.PostgresUser,
		DatabasePassword: cfg.PostgresPassword,
		DatabaseName:     cfg.PostgresDatabaseName,
		SSLMode:          postgres.SSLMode(cfg.PostgresSSLMode),
	}

	dbPool, err := postgres.Open(ctx, pgCfg)
	if err != nil {
		return fmt.Errorf("failed to open postgres connection: %w", err)
	}

	defer func() {
		logger.Infow("shutdown", "status", "closing DB connection")
		dbPool.Close()
	}()

	logger.Infow("DB migration [Starting...]")
	err = postgres.Migrate("file://scripts/db/migrations/", pgCfg)
	if err != nil {
		logger.Infow("postgres config", "config", pgCfg)
		logger.Errorw("unable to migrate postgres database", "error", err)
		return fmt.Errorf("failed to migrate postgres: %w", err)
	}
	logger.Infow("DB migration [Complete]")

	// -------------------------------------------------------------------------
	// Init service structs and handlers

	memberRepo := pgmembers.NewMembersRepository(logger, dbPool)
	membersUsecase := members.NewUsecase(memberRepo)

	classesRepo := pgclasses.NewClassesRepository(logger, dbPool)
	classesUsecase := classes.NewUsecase(classesRepo)

	bookingsRepo := pgbookings.NewBookingsRepository(logger, dbPool)
	bookingsUsecase := bookings.NewUsecase(bookingsRepo, membersUsecase, classesUsecase)

	pgProbe := postgres.NewProbe(dbPool)
	handlerCfg := handlers.Config{
		MembersUsecase: membersUsecase,
		ClassesUsecase: classesUsecase,
		BookingUsecase: bookingsUsecase,
		GinMode:        cfg.GinMode,
		Logger:         logger,
		PgProbe:        pgProbe,
	}
	handler, err := handlers.NewHandler(handlerCfg)
	if err != nil {
		return err
	}

	// -------------------------------------------------------------------------
	// Start http.Server
	server := http.Server{
		Addr:    cfg.Host,
		Handler: handler.API(),
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Infow("starting server", "host", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer logger.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("failed to stop server gracefully: %w", err)
		}
	}

	return nil
}
