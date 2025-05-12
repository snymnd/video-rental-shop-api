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
	ur    UsersRepository
	token TokenManager
}

type TokenManager interface {
	Generate(userID string, role int) (string, error)
}

func NewUsersUsecase(ur UsersRepository, token TokenManager) *UsersUsecase {
	return &UsersUsecase{
		ur:    ur,
		token: token,
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

func (uus *UsersUsecase) LoginUser(ctx context.Context, userLogin *entity.Login) error {
	// get user by email
	userLogin.Email = strings.ToLower(userLogin.Email)
	user, userErr := uus.ur.GetUserByEmail(ctx, userLogin.Email)
	if userErr != nil {
		if errors.Is(userErr, customerrors.ErrUserNotFound) {
			return customerrors.NewError(
				"email or password is invalid",
				userErr,
				customerrors.ItemNotExist,
			)
		}
	}
	userLogin.ID = user.ID
	userLogin.Name = user.Name
	userLogin.Role = user.Role

	// validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		return customerrors.NewError(
			"email or password is invalid",
			err,
			customerrors.CommonErr,
		)
	}

	// generate tokenStr
	tokenStr, tokenErr := uus.token.Generate(user.ID, user.Role)
	if tokenErr != nil {
		return customerrors.NewError("token generation fail", tokenErr, customerrors.CommonErr)
	}
	userLogin.Token = tokenStr

	return nil
}
