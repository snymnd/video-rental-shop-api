package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"
)

type RentalRepository struct {
	conn *sql.DB
}

func NewRentalRepository(conn *sql.DB) *RentalRepository {
	return &RentalRepository{conn}
}

func (pr *RentalRepository) Create(ctx context.Context, rentals entity.MultipleRentParams) error {
	query := `insert into rentals (video_id, payment_id, due_date)
				values `

	args := make([]any, len(rentals.VideosID)*3)
	count := 1
	for i, id := range rentals.VideosID {
		query += fmt.Sprintf("($%d, $%d, $%d),", count, count+1, count+2)
		args[i*3] = id
		args[i*3+1] = rentals.PaymentID
		args[i*3+2] = rentals.DueDate
		count += 3
	}
	query = query[:len(query)-1]

	var conn DBinf = pr.conn
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
