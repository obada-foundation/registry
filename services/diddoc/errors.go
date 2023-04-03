package diddoc

import (
	"errors"
)

var (
	// ErrDIDNotRegitered thows when requested DID is not present in registry
	ErrDIDNotRegitered = errors.New("DID not registered")

	// ErrDIDAlereadyRegistered thows when system tries to register DID that is already registered
	ErrDIDAlereadyRegistered = errors.New("DID was already registered")
)
