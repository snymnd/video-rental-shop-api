package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"
)

type RentalRepository struct {
	conn *sql.DB
}

func NewRentalRepository(conn *sql.DB) *RentalRepository {
	return &RentalRepository{conn}
}

func (rr *RentalRepository) Creates(ctx context.Context, rentals entity.MultipleRentParams) error {
	query := `insert into rentals (video_id, rental_payment_id, due_date, user_id)
				values `

	baseArgAmount := 4
	args := make([]any, len(rentals.VideosID)*baseArgAmount)
	count := 1
	for i, id := range rentals.VideosID {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d),", count, count+1, count+2, count+3)
		args[i*baseArgAmount] = id
		args[i*baseArgAmount+1] = rentals.PaymentID
		args[i*baseArgAmount+2] = rentals.DueDate
		args[i*baseArgAmount+3] = rentals.UserID
		count += baseArgAmount
	}
	query = query[:len(query)-1]

	var conn DBinf = rr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
	}

	_, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to create rental records",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}

func (rr *RentalRepository) UpdatesAddLatefee(ctx context.Context, rentalIDs []int, lateFeePaymentId int) error {
	query := `update rentals set latefee_payment_id = $1, updated_at = NOW()
				where id in (`

	args := make([]any, len(rentalIDs)+1)
	args[0] = lateFeePaymentId
	argsIdx := 1
	for _, id := range rentalIDs {
		query += fmt.Sprintf("$%d,", argsIdx+1)
		args[argsIdx] = id
		argsIdx += 1
	}
	query = query[:len(query)-1]
	query += ")"

	var conn DBinf = rr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
	}

	_, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to update rentals latefee payment id records",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}

func (rr *RentalRepository) FetchMultipleRentals(ctx context.Context, videosID []int, userID string, status constant.RentalStatus) (rentals entity.Rentals, err error) {
	query := `select id, video_id, user_id, rental_payment_id, latefee_payment_id, status, due_date, created_at, updated_at
				from rentals
				where user_id = $1 and video_id in (`

	args := make([]any, len(videosID)+1)
	args[0] = userID
	argsIdx := 1
	for _, id := range videosID {
		query += fmt.Sprintf("$%d,", argsIdx+1)
		args[argsIdx] = id
		argsIdx += 1
	}
	query = query[:len(query)-1]
	query += ")"

	if status != "" {
		query += fmt.Sprintf(" and status = $%d", argsIdx+1)
		args = append(args, status)
	}

	var conn DBinf = rr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
		query += " for update"
	}

	rows, rowsErr := conn.QueryContext(ctx, query, args...)
	if rowsErr != nil {
		return rentals, customerrors.NewError(
			"failed to fetch rentals",
			rowsErr,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	var rental entity.Rental
	for rows.Next() {
		if err = rows.Scan(
			&rental.ID,
			&rental.VideoID,
			&rental.UserID,
			&rental.RentalPaymentID,
			&rental.LateFeePaymentID,
			&rental.Status,
			&rental.DueDate,
			&rental.CreatedAt,
			&rental.UpdatedAt,
		); err != nil {
			return rentals, customerrors.NewError(
				"failed to scan rental data",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		rentals = append(rentals, rental)
	}

	if err = rows.Err(); err != nil {
		return rentals, customerrors.NewError(
			"failed to fetch rentals data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return rentals, nil
}

func (rr *RentalRepository) UpdatesRentalStatus(ctx context.Context, rentalIDs []int, status constant.RentalStatus) error {
	query := `update rentals set status = $1, updated_at = NOW()
				where id in (`

	args := make([]any, len(rentalIDs)+1)
	args[0] = status
	argsIdx := 1
	for _, id := range rentalIDs {
		query += fmt.Sprintf(`$%d,`, argsIdx+1)
		args[argsIdx] = id
		argsIdx += 1
	}
	query = query[:len(query)-1]
	query += ")"

	var conn DBinf = rr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
	}

	rows, rowsErr := conn.QueryContext(ctx, query, args...)
	if rowsErr != nil {
		return customerrors.NewError(
			"failed to update rentals status",
			rowsErr,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	return nil
}

func (rr *RentalRepository) UpdatesRentalStatusByPaymentID(ctx context.Context, paymentID int, status constant.RentalStatus) error {
	query := `update rentals set status = $1, updated_at = NOW()
				where rental_payment_id = $2`

	args := []any{status, paymentID}

	var conn DBinf = rr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
	}

	rows, rowsErr := conn.QueryContext(ctx, query, args...)
	if rowsErr != nil {
		return customerrors.NewError(
			"failed to update rentals status with specified id",
			rowsErr,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	return nil
}
