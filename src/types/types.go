package types

// RegisterDID describes payload for new DID registration
type RegisterDID struct {
	DID string `json:"did"`
}

// DIDDocument structure that expresses the DID Document
type DIDDocument struct {
	Context []string `json:"@context,omitempty" example:"https://www.w3.org/ns/did/v1"`
	ID      string   `json:"id,omitempty" example:"did:obada:123"`
}
