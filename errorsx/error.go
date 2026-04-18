package errorsx

import (
	"errors"
	"net/http"
)

const (
	CodeOK           = 0
	CodeBadRequest   = 40000
	CodeUnauthorized = 40100
	CodeForbidden    = 40300
	CodeNotFound     = 40400
	CodeInternal     = 50000
)

type Error struct {
	HTTPStatus int
	Code       int
	Msg        string
	Err        error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Msg
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(status int, code int, msg string) *Error {
	return &Error{
		HTTPStatus: status,
		Code:       code,
		Msg:        msg,
	}
}

func Wrap(err error, status int, code int, msg string) *Error {
	return &Error{
		HTTPStatus: status,
		Code:       code,
		Msg:        msg,
		Err:        err,
	}
}

func From(err error) *Error {
	if err == nil {
		return nil
	}
	var target *Error
	if errors.As(err, &target) {
		return target
	}
	return Wrap(err, http.StatusInternalServerError, CodeInternal, "internal server error")
}

var (
	ErrBadRequest   = New(http.StatusBadRequest, CodeBadRequest, "bad request")
	ErrUnauthorized = New(http.StatusUnauthorized, CodeUnauthorized, "unauthorized")
	ErrForbidden    = New(http.StatusForbidden, CodeForbidden, "forbidden")
	ErrNotFound     = New(http.StatusNotFound, CodeNotFound, "not found")
	ErrInternal     = New(http.StatusInternalServerError, CodeInternal, "internal server error")
)
