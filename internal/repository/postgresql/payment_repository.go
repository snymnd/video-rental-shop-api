package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"
)

type PaymentRepository struct {
	conn *sql.DB
}

func NewPaymentRepository(conn *sql.DB) *PaymentRepository {
	return &PaymentRepository{conn}
}

func (pr *PaymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	query := `insert into payments (user_id, total_price, expired_time)
				values ($1, $2, $3)
				returning id, status, created_at`

	var conn DBinf = pr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
	}

	if err := conn.QueryRowContext(ctx, query, payment.UserID, payment.TotalPrice, payment.ExpiredTime).
		Scan(
			&payment.ID,
			&payment.Status,
			&payment.CreatedAt,
		); err != nil {
		return customerrors.NewError(
			"cannot create payments",
			err,
			customerrors.DatabaseExecutionError,
		)

	}

	return nil
}

func (pr *PaymentRepository) GetPayment(ctx context.Context, paymentID int) (payment entity.Payment, err error) {
	query := `select id, user_id, total_price, method, expired_time, status, created_at
				from payments
				where id = $1`

	var conn DBinf = pr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
	}

	if err = conn.QueryRowContext(ctx, query, paymentID).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.TotalPrice,
		&payment.Method,
		&payment.ExpiredTime,
		&payment.Status,
		&payment.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return payment, customerrors.NewError(
				"payment is not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return payment, customerrors.NewError(
			"cannot get payments records",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return payment, nil
}

func (pr *PaymentRepository) UpdatePayment(ctx context.Context, updatePaymentParams entity.UpdatePaymentParams) error {
	query := `update payments set 
				status = coalesce($2, status),
				method = coalesce($3, method),
				total_price = coalesce($4, total_price),
				expired_time = coalesce($5, expired_time)
				where id = $1`

	var conn DBinf = pr.conn
	tx := ctx.Value(TxKey{TRANSACTION_KEY})
	if tx != nil {
		conn = tx.(*sql.Tx)
	}

	args := []any{updatePaymentParams.ID, updatePaymentParams.Status, updatePaymentParams.Method, updatePaymentParams.TotalPrice, updatePaymentParams.ExpiredTime}
	if _, err := conn.ExecContext(ctx, query, args...); err != nil {
		return customerrors.NewError(
			"cannot update payments",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}
