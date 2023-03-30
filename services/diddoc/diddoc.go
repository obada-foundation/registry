package diddoc

import (
	"context"

	immudb "github.com/codenotary/immudb/pkg/client"
	sdkdid "github.com/obada-foundation/sdkgo/did"
	"go.uber.org/zap"
)

// DIDDoc defines an API for work with DID documents
type DIDDoc interface {
	// Register registers a new DID document in the registry
	Register(ctx context.Context, did string) error

	// Get retrieves a DID document from the registry
	Get(ctx context.Context, did string) error
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
func (s Service) Register(ctx context.Context, did string) error {
	DID, err := sdkdid.FromString(did, nil)
	if err != nil {
		return err
	}

	_, err = s.db.Set(ctx, []byte(DID.String()), []byte("Foo"))
	if err != nil {
		return err
	}

	s.logger.Debugf("New DID registered: %q", DID)

	return nil
}

// Get implements DIDDoc Get
func (s Service) Get(ctx context.Context, did string) error {
	_, err := s.db.Get(ctx, []byte(did))
	if err != nil {
		return err
	}

	return nil
}
