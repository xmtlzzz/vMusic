package impl

import (
	"context"
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/xmtlzzz/vMusic/utils"
	"gorm.io/gorm"
)

type UserHandler struct {
	log zerolog.Logger
	db  *gorm.DB
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		log: zerolog.New(os.Stdout),
		db:  utils.NewDBConnector(),
	}
}

func (u *UserHandler) CreateUser(ctx context.Context, request *CreateUserRequest) error {
	if err := u.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	if err := gorm.G[User](u.db).Create(ctx, &User{Username: request.Username,
		Mail:      request.Mail,
		Telephone: request.Telephone,
		Password:  request.Password}); err != nil {
		return err
	}
	return nil
}

func (u *UserHandler) DeleteUser(ctx context.Context, request *DeleteUserRequest) (int, error) {
	if request.Username == "" {
		return 0, errors.New("username is required")
	}
	dL, err := gorm.G[User](u.db).Where("user = ?", request.Username).Delete(ctx)
	if err != nil {
		return 0, err
	}
	return dL, nil
}

func (u *UserHandler) UserIsAvailable(ctx context.Context, request *UserLoginRequest) (*User, error) {
	user, err := gorm.G[User](u.db).Where("username = ?", request.Username).First(ctx)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Password != request.Password {
		return nil, errors.New("username or password is wrong")
	}
	return &user, nil
}
