package impl

import "context"

type TokenInt interface {
	RegistryToken(ctx context.Context, request *RegistryTokenRequest) (*Token, error)
	DeleteToken(ctx context.Context, request *DeleteTokenRequest) error
}
