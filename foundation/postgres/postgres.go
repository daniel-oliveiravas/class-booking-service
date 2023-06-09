package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	minConnectionsDefault    = 1
	maxConnectionsDefault    = 3
	healthCheckPeriodDefault = time.Second * 30
	sslModeDefault           = DisableSSLMode
)

type Config struct {
	Host              string
	Port              int
	DatabaseUser      string
	DatabasePassword  string
	DatabaseName      string
	MaxConnections    int32
	MinConnections    int32
	HealthCheckPeriod time.Duration
	SSLMode           SSLMode
	SearchPath        string
}

type SSLMode string

const (
	DisableSSLMode SSLMode = "disable"
)

func (c *Config) toURL() string {
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s/%s",
		url.QueryEscape(c.DatabaseUser),
		url.QueryEscape(c.DatabasePassword),
		net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		url.QueryEscape(c.DatabaseName))

	var options []string
	if c.SSLMode != "" {
		options = append(options, fmt.Sprintf("sslmode=%s", sslModeDefault))
	}
	if c.SearchPath != "" {
		options = append(options, fmt.Sprintf("search_path=%s", c.SearchPath))
	}

	for idx, option := range options {
		if idx == 0 {
			postgresURL = fmt.Sprintf("%s?%s", postgresURL, option)
		}
		postgresURL = fmt.Sprintf("%s&%s", postgresURL, option)
	}

	return postgresURL
}

func Open(ctx context.Context, conf Config) (*pgxpool.Pool, error) {
	pgDatabaseURL := conf.toURL()

	pgConfig, err := pgxpool.ParseConfig(pgDatabaseURL)
	if err != nil {
		return nil, err
	}

	pgConfig.MinConns = minConnectionsDefault
	pgConfig.MaxConns = maxConnectionsDefault

	if conf.MinConnections != 0 {
		pgConfig.MinConns = conf.MinConnections
	}

	if conf.MaxConnections != 0 {
		pgConfig.MaxConns = conf.MaxConnections
	}

	if conf.MinConnections > conf.MaxConnections {
		conf.MinConnections = conf.MaxConnections
	}

	pgConfig.HealthCheckPeriod = healthCheckPeriodDefault
	if conf.HealthCheckPeriod != 0 {
		pgConfig.HealthCheckPeriod = conf.HealthCheckPeriod
	}

	var pgPool *pgxpool.Pool
	pgPool, err = pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return nil, err
	}

	return pgPool, nil
}

func Migrate(migrationSource string, conf Config) error {
	m, err := migrate.New(migrationSource, conf.toURL())
	if err != nil {
		return fmt.Errorf("could not migrate database %s. err: %w", conf.DatabaseName, err)
	}

	defer m.Close()

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		version, dirty, _ := m.Version()
		return fmt.Errorf("failed to run database migration (current version: %d, dirty: %v, %w)", version, dirty, err)
	}

	return nil
}

func DropAndCreateSchema(ctx context.Context, pool *pgxpool.Pool, schema string) error {
	_, err := pool.Exec(ctx, fmt.Sprintf("DROP schema IF EXISTS %s cascade;", schema))
	if err != nil {
		return fmt.Errorf("could not drop schema %s: %w", schema, err)
	}
	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE schema %s;", schema))
	if err != nil {
		return fmt.Errorf("could not create schema %s: %w", schema, err)
	}
	return nil
}
