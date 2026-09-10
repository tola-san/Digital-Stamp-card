package domain

import "errors"

var (
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrPhoneAlreadyExists = errors.New("phone number is already registered")
	ErrInvalidName        = errors.New("name must contain between 2 and 100 characters")
	ErrInvalidPhone       = errors.New("phone number must use international format")
	ErrInvalidSession     = errors.New("customer session is invalid")
	ErrInvalidCursor      = errors.New("transaction cursor is invalid")
	ErrInvalidPageLimit   = errors.New("transaction limit must be between 1 and 100")
)
