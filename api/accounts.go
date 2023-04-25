package api

import (
	"context"

	pb "github.com/obada-foundation/registry/api/pb/v1/account"
)

// RegisterAccount registers a new account in a registry
func (s GRPCServer) RegisterAccount(ctx context.Context, msg *pb.RegisterAccountRequest) (*pb.RegisterAccountResponse, error) {

	if err := s.AccountService.Register(ctx, msg.GetPubkey()); err != nil {
		return nil, err
	}

	return &pb.RegisterAccountResponse{}, nil
}

// GetPublicKey returns a public key for registry account
func (s GRPCServer) GetPublicKey(ctx context.Context, msg *pb.GetPublicKeyRequest) (*pb.GetPublicKeyResponse, error) {
	pubKey, err := s.AccountService.GetPublicKey(ctx, msg.GetAddress())
	if err != nil {
		return nil, err
	}

	return &pb.GetPublicKeyResponse{
		Pubkey: pubKey,
	}, nil
}
