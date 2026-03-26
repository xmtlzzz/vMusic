package impl

import "context"

type UserInter interface {
	CreateUser(ctx context.Context, request *CreateUserRequest) error
	DeleteUser(ctx context.Context, request *DeleteUserRequest) (int, error)
	UserIsAvailable(ctx context.Context, request *UserLoginRequest) (*User, error)
}
