package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ClassesRepository struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

func NewClassesRepository(logger *zap.SugaredLogger, db *pgxpool.Pool) *ClassesRepository {
	return &ClassesRepository{
		logger: logger,
		db:     db,
	}
}

func (r *ClassesRepository) Add(ctx context.Context, class classes.Class) (classes.Class, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return classes.Class{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	insertClass := `INSERT INTO classes (id, name, start_date, end_date, capacity) 
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, name, created_at, updated_at, start_date, end_date, capacity`
	row := tx.QueryRow(ctx, insertClass, class.ID, class.Name, class.StartDate, class.EndDate, class.Capacity)

	var storesClass classes.Class
	err = row.Scan(&storesClass.ID, &storesClass.Name, &storesClass.CreatedAt, &storesClass.UpdatedAt, &storesClass.StartDate, &storesClass.EndDate, &storesClass.Capacity)
	if err != nil {
		return classes.Class{}, fmt.Errorf("failed to scan classes row to classes.Class: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return classes.Class{}, fmt.Errorf("failed to get booking by ID: %w", err)
	}
	return storesClass, nil
}

func (r *ClassesRepository) GetByID(ctx context.Context, classID string) (classes.Class, error) {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return classes.Class{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	class, err := r.getByIdTxn(ctx, txn, classID)
	if err != nil {
		return classes.Class{}, err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return classes.Class{}, fmt.Errorf("failed to get class by ID: %w", err)
	}

	return class, nil
}

func (r *ClassesRepository) getByIdTxn(ctx context.Context, txn pgx.Tx, classID string) (classes.Class, error) {
	query := `SELECT id, created_at, updated_at, name, start_date, end_date, capacity FROM classes WHERE id = $1;`

	row := txn.QueryRow(ctx, query, classID)
	var class classes.Class
	err := row.Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt, &class.Name, &class.StartDate, &class.EndDate, &class.Capacity)
	if err != nil {
		return classes.Class{}, fmt.Errorf("failed to scan classes row to classes.Class: %w", err)
	}
	return class, nil
}

func (r *ClassesRepository) IsNotFoundErr(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func (r *ClassesRepository) Update(ctx context.Context, classID string, updateClass classes.UpdateClass) (classes.Class, error) {
	values := make([]interface{}, 0)
	columns := make([]string, 0)

	if updateClass.Name != nil {
		values = append(values, *updateClass.Name)
		columns = append(columns, "name")
	}

	if updateClass.StartDate != nil {
		values = append(values, *updateClass.StartDate)
		columns = append(columns, "start_date")
	}

	if updateClass.EndDate != nil {
		values = append(values, *updateClass.EndDate)
		columns = append(columns, "end_date")
	}

	if updateClass.Capability != nil {
		values = append(values, *updateClass.Capability)
		columns = append(columns, "capacity")
	}

	if len(values) == 0 {
		return classes.Class{}, nil
	}

	updateStatements := make([]string, 0)
	for index := range columns {
		updateStatements = append(updateStatements, fmt.Sprintf("%s = $%d", columns[index], index+1))
	}

	values = append(values, classID)
	var statement = "UPDATE classes SET " + strings.Join(updateStatements, ", ") + fmt.Sprintf(" WHERE id = $%d", len(values)) +
		" RETURNING id, created_at, updated_at, name, start_date, end_date, capacity"

	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return classes.Class{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer txn.Rollback(ctx)

	row := r.db.QueryRow(ctx, statement, values...)
	var class classes.Class
	err = row.Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt, &class.Name, &class.StartDate, &class.EndDate, &class.Capacity)
	if err != nil {
		return classes.Class{}, fmt.Errorf("failed to scan classes row to classes.Class: %w", err)
	}

	if err := txn.Commit(ctx); err != nil {
		return classes.Class{}, fmt.Errorf("failed to commit transaction. :%w", err)
	}

	return class, nil
}

func (r *ClassesRepository) Delete(ctx context.Context, classID string) error {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer txn.Rollback(ctx)

	statement := `DELETE FROM classes WHERE id = $1`
	_, err = txn.Exec(ctx, statement, classID)
	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction. :%w", err)
	}

	return nil
}

func (r *ClassesRepository) List(ctx context.Context, limit int, offset int) ([]classes.Class, error) {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `SELECT id, created_at, updated_at, name, start_date, end_date, capacity
				FROM classes
			  LIMIT $1 OFFSET $2;`

	rows, err := txn.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query classes: %w", err)
	}

	allClasses := make([]classes.Class, 0)
	for rows.Next() {
		var class classes.Class
		err := rows.Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt, &class.Name, &class.StartDate, &class.EndDate, &class.Capacity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan classes row to classes.Class: %w", err)
		}

		allClasses = append(allClasses, class)
	}

	if err := txn.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction. :%w", err)
	}

	return allClasses, nil
}
