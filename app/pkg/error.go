package pkg

import (
	"event-booking-api/app/constant"
	"fmt"
)

type CustomError struct {
	Type constant.ResponseStatus
	Msg  string
	Err  error
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return fmt.Sprint(e.Msg)
}

func NewCustomError(typ constant.ResponseStatus, msg string, err error) *CustomError {
	return &CustomError{Type: typ, Msg: msg, Err: err}
}

func NewNotFoundError(msg string, err error) *CustomError {
	return NewCustomError(constant.DataNotFound, msg, err)
}

func NewConflictError(msg string, err error) *CustomError {
	return NewCustomError(constant.Conflict, msg, err)
}

func NewUnauthorizedError(msg string, err error) *CustomError {
	return NewCustomError(constant.Unauthorized, msg, err)
}
