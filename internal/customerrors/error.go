package customerrors

import (
	"errors"
	"net/http"
	"vrs-api/internal/dto"
)

type CustomError struct {
	ErrorMessage string
	Details      any
	ErrorLog     error
	ErrorCode    int
}

func (err CustomError) Error() string {
	return err.ErrorMessage
}

func NewError(message string, errorLog error, errorCode int, details ...[]dto.DetailsError) *CustomError {

	var detail any
	if len(details) > 0 {
		detail = details[0]
	}

	return &CustomError{
		ErrorMessage: message,
		ErrorLog:     errorLog,
		ErrorCode:    errorCode,
		Details:      detail,
	}
}

func (err CustomError) GetHTTPErrorCode() int {
	switch err.ErrorCode {
	case ItemAlreadyExist:
		return http.StatusBadRequest
	case ItemNotExist:
		return http.StatusNotFound
	case InvalidAction:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Unauthenticate:
		return http.StatusUnauthorized
	case DatabaseExecutionError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

const (
	ItemAlreadyExist       = 40400
	InvalidAction          = 41400
	Unauthenticate         = 40401
	Unauthorized           = 40402
	ItemNotExist           = 40404
	CommonErr              = 50500
	DatabaseExecutionError = 51500
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrFailedGetAuthPayload = NewError(
		"user credential identification failed",
		errors.New("failed to get authorization payload from context"),
		CommonErr,
	)
	ErrUserIDNotFoundInAuthPayload = NewError(
		"user credential identification failed",
		errors.New("user id cannot be found on auth payload subject"),
		CommonErr,
	)
)
