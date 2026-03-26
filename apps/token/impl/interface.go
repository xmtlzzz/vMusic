package impl

import (
	"context"

	user_impl "github.com/xmtlzzz/vMusic/apps/user/impl"
)

type TokenInt interface {
	RegistryToken(ctx context.Context, request *RegistryTokenRequest) (*Token, error)
	QueryUserByAccessToken(ctx context.Context, request *QueryUserByTokenRequest) (*user_impl.User, error)
	DeleteToken(ctx context.Context, request *DeleteTokenRequest) error
}
