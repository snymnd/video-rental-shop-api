package postgresql

import (
	"context"
	"database/sql"
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
		return err
	}

	return nil
}
