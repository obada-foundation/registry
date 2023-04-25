package cmd

// ClientCommand with command line flags and env vars
type ClientCommand struct {
	Account AccountCommand `command:"account" description:"Working with registry accounts(public keys)"`

	DIDDoc DIDDocCommand `command:"diddoc" description:"Working with DID documents"`

	CommonOpts
}

// AccountCommand wraps functionality for working with accounts
type AccountCommand struct {
	Register AccountRegisterCommand `command:"register" description:"Registers a new public key/account in registry"`

	PublicKey AccountPublicKeyCommand `command:"public-key" description:"Fetches a public key by address"`
}

// AccountRegisterCommand Registers a new public key/account in registry"
type AccountRegisterCommand struct {
}

// AccountPublicKeyCommand fetches a public key by address
type AccountPublicKeyCommand struct {
}

// DIDDocCommand wraps functionality for working with DID documents
type DIDDocCommand struct {
}

// Execute is the entry point for "client" command, called by flag parser
func (c *ClientCommand) Execute(_ []string) error {
	return nil
}
