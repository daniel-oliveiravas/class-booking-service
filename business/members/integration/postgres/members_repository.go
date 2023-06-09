package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type MembersRepository struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

func NewMembersRepository(logger *zap.SugaredLogger, db *pgxpool.Pool) *MembersRepository {
	return &MembersRepository{
		logger: logger,
		db:     db,
	}
}

func (r *MembersRepository) AddMember(ctx context.Context, member members.Member) (members.Member, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return members.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	insertMember := `INSERT INTO members (id, name) 
				VALUES ($1, $2)
				RETURNING id, name, created_at, updated_at`
	row := tx.QueryRow(ctx, insertMember, member.ID, member.Name)

	var storedMember members.Member
	err = row.Scan(&storedMember.ID, &storedMember.Name, &storedMember.CreatedAt, &storedMember.UpdatedAt)
	if err != nil {
		return members.Member{}, fmt.Errorf("failed to scan members row to members.Member: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return members.Member{}, fmt.Errorf("failed to get booking by ID: %w", err)
	}
	return storedMember, nil
}

func (r *MembersRepository) GetByID(ctx context.Context, memberID string) (members.Member, error) {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return members.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	member, err := r.getByIdTxn(ctx, txn, memberID)
	if err != nil {
		return members.Member{}, err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return members.Member{}, fmt.Errorf("failed to get member by ID: %w", err)
	}

	return member, nil
}

func (r *MembersRepository) getByIdTxn(ctx context.Context, txn pgx.Tx, memberID string) (members.Member, error) {
	query := `SELECT id, created_at, updated_at, name FROM members WHERE id = $1;`

	row := txn.QueryRow(ctx, query, memberID)
	var member members.Member
	err := row.Scan(&member.ID, &member.CreatedAt, &member.UpdatedAt, &member.Name)
	if err != nil {
		return members.Member{}, fmt.Errorf("failed to scan members row to members.Member: %w", err)
	}
	return member, nil
}

func (r *MembersRepository) IsNotFoundErr(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func (r *MembersRepository) UpdateMember(ctx context.Context, memberID string, updateMember members.UpdateMember) (members.Member, error) {
	var member members.Member
	if updateMember.Name != nil {
		txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
			AccessMode: pgx.ReadWrite,
		})
		if err != nil {
			return members.Member{}, fmt.Errorf("failed to begin transaction: %w", err)
		}

		defer txn.Rollback(ctx)

		statement := `UPDATE members SET name = $1, updated_at = now() WHERE id = $2
						RETURNING id, name, created_at, updated_at`
		row := txn.QueryRow(ctx, statement, updateMember.Name, memberID)

		err = row.Scan(&member.ID, &member.Name, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			return members.Member{}, fmt.Errorf("failed to update member: %w", err)
		}

		if err := txn.Commit(ctx); err != nil {
			return members.Member{}, fmt.Errorf("failed to commit transaction. :%w", err)
		}
	}

	return member, nil
}

func (r *MembersRepository) DeleteMember(ctx context.Context, memberID string) error {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer txn.Rollback(ctx)

	statement := `DELETE FROM members WHERE id = $1`
	_, err = txn.Exec(ctx, statement, memberID)
	if err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}

	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction. :%w", err)
	}

	return nil
}

func (r *MembersRepository) ListMembers(ctx context.Context, limit int, offset int) ([]members.Member, error) {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `SELECT id, created_at, updated_at, name
				FROM members m
			  LIMIT $1 OFFSET $2;`

	rows, err := txn.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}

	allMembers := make([]members.Member, 0)
	for rows.Next() {
		var member members.Member
		err := rows.Scan(&member.ID, &member.CreatedAt, &member.UpdatedAt, &member.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan members row to members.Member: %w", err)
		}

		allMembers = append(allMembers, member)
	}

	if err := txn.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction. :%w", err)
	}

	return allMembers, nil
}
