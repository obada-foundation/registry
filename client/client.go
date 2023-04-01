package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/obada-foundation/registry/types"
)

// Client is allows to query the registry
type Client interface {
	// Register a DID with the registry
	Register(DID string) error

	// Get a DID document from the registry
	Get(DID string) (types.DIDDocument, error)
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

// Get implements Client.Get
func (c *HTTPClient) Get(did string) (types.DIDDocument, error) {
	var doc types.DIDDocument

	resp, err := c.h.Get(c.url + "/api/v1.0/" + did)
	if err != nil {
		return doc, err
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&doc); err != nil {
		return doc, err
	}
	_ = resp.Body.Close()

	return doc, nil
}
