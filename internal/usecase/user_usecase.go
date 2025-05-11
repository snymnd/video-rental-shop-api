package usecase

import (
	"context"
	"errors"
	"strings"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

type UsersRepository interface {
	Create(ctx context.Context, user *entity.Users) error
	CheckIsEmailExist(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.Users, error)
}

type UsersUsecase struct {
	ur UsersRepository
}

type Token interface {
	Generate(userID string) (string, error)
}

func NewUsersUsecase(ur UsersRepository) *UsersUsecase {
	return &UsersUsecase{
		ur: ur,
	}
}

func (uus *UsersUsecase) RegisterUser(ctx context.Context, user *entity.Users) error {
	// check is email exist
	email := strings.ToLower(user.Email)
	isExist, err := uus.ur.CheckIsEmailExist(ctx, email)
	if err != nil {
		return err
	}

	if isExist {
		return customerrors.NewError(
			"email is already exist",
			errors.New("email is already is exist in users data"),
			customerrors.ItemAlreadyExist,
		)
	}

	// hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Email = strings.ToLower(user.Email)
	if err := uus.ur.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

