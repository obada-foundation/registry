package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/types"
	"github.com/obada-foundation/sdkgo/asset"
)

// Client is allows to query the registry
type Client interface {
	// Register a DID with the registry
	Register(newDID types.RegisterDID) error

	// Get a DID document from the registry
	Get(DID string) (types.DIDDocument, error)

	// GetMetadataHistory returns the history of changes of asset data
	GetMetadataHistory(DID string) (asset.DataArrayVersions, error)

	// SaveMetadata saves the asset metadata to the regostry
	SaveMetadata(DID string, md types.SaveMetadata) error
}

// HTTPClient is a client for the registry API
type HTTPClient struct {
	h   *http.Client
	url string
}

// NewHTTPClient creates a new instance of HTTPClient
func NewHTTPClient(registryURL string) *HTTPClient {
	h := &http.Client{}

	return &HTTPClient{
		h:   h,
		url: registryURL,
	}
}

// Register implements Client.Register
func (c *HTTPClient) Register(newDID types.RegisterDID) error {
	b, err := json.Marshal(newDID)
	if err != nil {
		return err
	}

	resp, err := c.h.Post(c.url+"/api/v1.0/register", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to register DID: %w", err)
		}
		_ = resp.Body.Close()

		return fmt.Errorf("failed to register DID: %q", string(b))
	}

	return nil
}

// SaveMetadata implements Client.SaveMetadata
func (c *HTTPClient) SaveMetadata(did string, md types.SaveMetadata) error {
	b, err := json.Marshal(md)
	if err != nil {
		return err
	}

	resp, err := c.h.Post(c.url+"/api/v1.0/"+did+"/metadata", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save metadata: %w", err)
		}
		_ = resp.Body.Close()

		return fmt.Errorf("failed to save metadata: %q", string(b))
	}

	return nil
}

// GetMetadataHistory implements Client.GetMetadataHistory
func (c *HTTPClient) GetMetadataHistory(did string) (asset.DataArrayVersions, error) {
	var history asset.DataArrayVersions

	resp, err := c.h.Get(c.url + "/api/v1.0/" + did + "/metadata-history")
	if err != nil {
		return history, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return history, diddoc.ErrDIDNotRegistered
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&history); err != nil {
		return history, err
	}
	_ = resp.Body.Close()

	return history, nil
}

// Get implements Client.Get
func (c *HTTPClient) Get(did string) (types.DIDDocument, error) {
	var doc types.DIDDocument

	resp, err := c.h.Get(c.url + "/api/v1.0/" + did)
	if err != nil {
		return doc, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return doc, diddoc.ErrDIDNotRegistered
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&doc); err != nil {
		return doc, err
	}
	_ = resp.Body.Close()

	return doc, nil
}
