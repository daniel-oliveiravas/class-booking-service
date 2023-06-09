package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Probe struct {
	pool *pgxpool.Pool
}

func NewProbe(pool *pgxpool.Pool) *Probe {
	return &Probe{pool: pool}
}

func (db *Probe) Check(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
