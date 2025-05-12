package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"vrs-api/internal/customerrors"
)

type RBACRepository struct {
	conn *sql.DB
}

func NewRBACRepository(conn *sql.DB) *RBACRepository {
	return &RBACRepository{conn}
}
func (rr *RBACRepository) HasAccess(ctx context.Context, role int, permission int, resource int) (bool, error) {
	query := `select id 
				from rbac
				where role_id = $1
				and permission_id = $2
				and resource_id = $3;`

	if err := rr.conn.QueryRowContext(ctx, query, role, permission, resource).Scan(new(int)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, customerrors.NewError(
			"cannot check access unique validation",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return true, nil
}
