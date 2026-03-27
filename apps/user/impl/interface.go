package impl

import (
	"context"

	music "github.com/xmtlzzz/vMusic/apps/music/impl"
)

type UserInter interface {
	CreateUser(ctx context.Context, request *CreateUserRequest) error
	DeleteUser(ctx context.Context, request *DeleteUserRequest) (int, error)
	UserIsAvailable(ctx context.Context, request *UserLoginRequest) (*User, error)
	GetUserByID(ctx context.Context, userID int) (*User, error)
	ToggleLike(ctx context.Context, userID int, track *music.Objector) (*ToggleLikeResponse, error)
}
