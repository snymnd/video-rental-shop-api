package customerrors

import (
	"errors"
	"net/http"
)

type CustomError struct {
	ErrorMessage string
	ErrorLog     error
	ErrorCode    int
}

func (err CustomError) Error() string {
	return err.ErrorMessage
}

// error builder for a known error
func NewError(message string, errorLog error, errorCode int) *CustomError {
	return &CustomError{
		ErrorMessage: message,
		ErrorLog:     errorLog,
		ErrorCode:    errorCode,
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
	ErrUserNotFound = errors.New("user not found")
)
