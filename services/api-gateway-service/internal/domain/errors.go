package domain

import (
	"errors"
	"fmt"
)

var (
	ErrSatelliteNotFound   = errors.New("satellite not found")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrUpstreamUnavailable = errors.New("upstream unavailable")
	ErrUpstreamTimeout     = errors.New("upstream timeout")
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(format string, args ...any) error {
	return &ValidationError{
		Message: fmt.Sprintf(format, args...),
	}
}
