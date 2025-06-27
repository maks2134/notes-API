package util

import "errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrValidationFailed = errors.New("validation failed")
	ErrNotFound         = errors.New("not found")
	ErrInternalServer   = errors.New("internal server error")
	ErrBadRequest       = errors.New("bad request")
)
