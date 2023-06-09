package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/daniel-oliveiravas/class-booking-service/business/bookings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BookingsRepository struct {
	logger *zap.SugaredLogger
	db     *pgxpool.Pool
}

func NewBookingsRepository(logger *zap.SugaredLogger, db *pgxpool.Pool) *BookingsRepository {
	return &BookingsRepository{
		logger: logger,
		db:     db,
	}
}

func (r *BookingsRepository) BookClass(ctx context.Context, booking bookings.Booking) (bookings.Booking, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return bookings.Booking{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	insertBooking := `INSERT INTO bookings (id, member_id, class_id, class_date)
				VALUES ($1, $2, $3, $4)
				RETURNING id, booked_at, updated_at, member_id, class_id, class_date`
	row := tx.QueryRow(ctx, insertBooking, booking.ID, booking.MemberID, booking.ClassID, booking.ClassDate)

	var storedBooking bookings.Booking
	err = row.Scan(&storedBooking.ID, &storedBooking.BookedAt, &storedBooking.UpdatedAt, &storedBooking.MemberID, &storedBooking.ClassID, &storedBooking.ClassDate)
	if err != nil {
		return bookings.Booking{}, fmt.Errorf("failed to scan bookings row to bookings.Booking: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return bookings.Booking{}, fmt.Errorf("failed to get booking by ID: %w", err)
	}

	return storedBooking, nil
}

func (r *BookingsRepository) GetByID(ctx context.Context, bookingID string) (bookings.Booking, error) {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return bookings.Booking{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	booking, err := r.getByIdTxn(ctx, txn, bookingID)
	if err != nil {
		return bookings.Booking{}, err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return bookings.Booking{}, fmt.Errorf("failed to get booking by ID: %w", err)
	}

	return booking, nil
}

func (r *BookingsRepository) getByIdTxn(ctx context.Context, txn pgx.Tx, bookingID string) (bookings.Booking, error) {
	query := `SELECT id, booked_at, updated_at, member_id, class_id, class_date FROM bookings WHERE id = $1;`

	row := txn.QueryRow(ctx, query, bookingID)
	var booking bookings.Booking
	err := row.Scan(&booking.ID, &booking.BookedAt, &booking.UpdatedAt, &booking.MemberID, &booking.ClassID, &booking.ClassDate)
	if err != nil {
		return bookings.Booking{}, fmt.Errorf("failed to scan bookings row to bookings.Booking: %w", err)
	}
	return booking, nil
}

func (r *BookingsRepository) IsNotFoundErr(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func (r *BookingsRepository) DeleteBooking(ctx context.Context, bookingID string) error {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer txn.Rollback(ctx)

	statement := `DELETE FROM bookings WHERE id = $1`
	_, err = txn.Exec(ctx, statement, bookingID)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction. :%w", err)
	}

	return nil
}

func (r *BookingsRepository) ListBookings(ctx context.Context, limit int, offset int) ([]bookings.Booking, error) {
	txn, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `SELECT id, booked_at, updated_at, member_id, class_id, class_date
				FROM bookings m
			  LIMIT $1 OFFSET $2;`

	rows, err := txn.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}

	allbookings := make([]bookings.Booking, 0)
	for rows.Next() {
		var booking bookings.Booking
		err := rows.Scan(&booking.ID, &booking.BookedAt, &booking.UpdatedAt, &booking.MemberID, &booking.ClassID, &booking.ClassDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bookings row to bookings.Booking: %w", err)
		}

		allbookings = append(allbookings, booking)
	}

	if err := txn.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction. :%w", err)
	}

	return allbookings, nil
}
