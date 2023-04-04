package diddoc

import (
	"errors"
)

var (
	// ErrDIDNotRegistered thows when requested DID is not present in registry
	ErrDIDNotRegistered = errors.New("DID not registered")

	// ErrDIDAlereadyRegistered thows when system tries to register DID that is already registered
	ErrDIDAlereadyRegistered = errors.New("DID was already registered")
)
