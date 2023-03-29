package diddoc

import (
	immudb "github.com/codenotary/immudb/pkg/client"
	sdkdid "github.com/obada-foundation/sdkgo/did"
	"go.uber.org/zap"
)

// DIDDoc defines an API for work with DID documents
type DIDDoc interface {
	// Register registers a new DID document in the registry
	Register(did string) error
}

// Service implements DIDDoc
type Service struct {
	db     immudb.ImmuClient
	logger *zap.SugaredLogger
}

// NewService creates a new instance of DIDDoc service
func NewService(db immudb.ImmuClient, logger *zap.SugaredLogger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

// Register implements DIDDoc Register
func (s Service) Register(did string) error {
	DID, err := sdkdid.FromString(did, nil)
	if err != nil {
		return err
	}

	s.logger.Debugf("New DID registered", DID)

	return nil
}
