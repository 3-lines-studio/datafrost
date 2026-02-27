package usecase

import "errors"

var (
	ErrConnectionNotFound = errors.New("connection not found")
	ErrQueryNotFound      = errors.New("query not found")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrNameRequired       = errors.New("name is required")
	ErrTypeRequired       = errors.New("type is required")
	ErrQueryRequired      = errors.New("query is required")
)
