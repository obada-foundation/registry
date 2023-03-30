package diddoc

import (
	"bytes"
	"context"
	"encoding/gob"

	immudb "github.com/codenotary/immudb/pkg/client"
	"github.com/obada-foundation/registry/system/encoder"
	"github.com/obada-foundation/registry/types"
	sdkdid "github.com/obada-foundation/sdkgo/did"
	"go.uber.org/zap"
)

// DIDDoc defines an API for work with DID documents
type DIDDoc interface {
	// Register registers a new DID document in the registry
	Register(ctx context.Context, did string) error

	// Get retrieves a DID document from the registry
	Get(ctx context.Context, did string) (types.DIDDocument, error)
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

	gobData, err := encoder.DataEncode(types.DIDDocument{
		ID: DID.String(),
	})
	if err != nil {
		return err
	}

	_, err = s.db.Set(ctx, []byte(DID.String()), gobData)
	if err != nil {
		return err
	}

	s.logger.Debugf("New DID registered: %q", DID)

	return nil
}

// Get implements DIDDoc Get
func (s Service) Get(ctx context.Context, did string) (types.DIDDocument, error) {
	var DIDDoc types.DIDDocument

	entry, err := s.db.Get(ctx, []byte(did))
	if err != nil {
		return DIDDoc, err
	}

	dec := gob.NewDecoder(bytes.NewBuffer(entry.Value))
	if err := dec.Decode(&DIDDoc); err != nil {
		return DIDDoc, err
	}

	return DIDDoc, nil
}
