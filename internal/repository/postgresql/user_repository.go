package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"
)

type UserRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn}
}

func (ur *UserRepository) Create(ctx context.Context, user *entity.Users) error {
	query := `insert into users (name, email, password) 
				values ($1, $2, $3) returning id, role, created_at, updated_at;`

	if err := ur.conn.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).
		Scan(
			&user.ID,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
		return customerrors.NewError(
			"cannot create user data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}

func (br *UserRepository) CheckIsEmailExist(ctx context.Context, email string) (bool, error) {
	query := `select email 
				from users
				where email = $1`

	if err := br.conn.QueryRowContext(ctx, query, email).Scan(new(string)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, customerrors.NewError(
			"cannot check email unique validation",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return true, nil
}

func (br *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.Users, error) {
	query := `select id, name, email, password, role, created_at, updated_at 
				from users
				where LOWER(email) = $1`
	var user entity.Users

	if err := br.conn.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ErrUserNotFound
		}

		return nil, customerrors.NewError(
			"cannot get user with email",
			errors.New("cannot get users with email in users data"),
			customerrors.ItemNotExist,
		)
	}

	return &user, nil
}

func (ur *UserRepository) CheckIsUserExist(ctx context.Context, id string) error {
	query := `select id 
				from users
				where id = $1`

	var a string
	if err := ur.conn.QueryRowContext(ctx, query, id).Scan(&a); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"user not found",
				err,
				customerrors.ItemNotExist,
			)
		}

		return customerrors.NewError(
			"cannot check user existance",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}
