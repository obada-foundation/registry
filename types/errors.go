package types

import (
	"errors"
)

var (
	ErrUnauthorized            = errors.New("unauthorized")
	ErrUnauthorizedNoSignature = errors.New("unauthorized, no signature")
)
