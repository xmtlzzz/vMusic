package impl

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog"
	music "github.com/xmtlzzz/vMusic/apps/music/impl"
	"github.com/xmtlzzz/vMusic/utils"
	"gorm.io/gorm"
)

type UserHandler struct {
	log zerolog.Logger
	db  *gorm.DB
}

func NewUserHandler() *UserHandler {
	handler := &UserHandler{
		log: zerolog.New(os.Stdout),
		db:  utils.NewDBConnector(),
	}
	if err := handler.ensureSchema(); err != nil {
		handler.log.Error().Err(err).Msg("failed to migrate user schema")
	}
	return handler
}

func (u *UserHandler) CreateUser(ctx context.Context, request *CreateUserRequest) error {
	if err := u.ensureSchema(); err != nil {
		return err
	}

	username := strings.TrimSpace(request.UserInfo.Username)
	password := strings.TrimSpace(request.UserInfo.Password)
	if username == "" || password == "" {
		return errors.New("username and password are required")
	}

	user := &User{
		Username:  username,
		Mail:      strings.TrimSpace(request.UserInfo.Mail),
		Telephone: strings.TrimSpace(request.UserInfo.Telephone),
		Password:  password,
		LikeList:  LikeList{},
		IsAdmin:   request.UserInfo.IsAdmin,
	}

	if err := gorm.G[User](u.db).Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("user is created")
		}
		return err
	}
	return nil
}

func (u *UserHandler) DeleteUser(ctx context.Context, request *DeleteUserRequest) (int, error) {
	if request.Username == "" {
		return 0, errors.New("username is required")
	}
	if err := u.ensureSchema(); err != nil {
		return 0, err
	}
	dL, err := gorm.G[User](u.db).Where("username = ?", request.Username).Delete(ctx)
	if err != nil {
		return 0, err
	}
	return dL, nil
}

func (u *UserHandler) UserIsAvailable(ctx context.Context, request *UserLoginRequest) (*User, error) {
	if err := u.ensureSchema(); err != nil {
		return nil, err
	}
	user, err := gorm.G[User](u.db).Where("username = ?", request.Username).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if user.Password != request.Password {
		return nil, errors.New("username or password is wrong")
	}
	return user.Sanitized(), nil
}

func (u *UserHandler) GetUserByID(ctx context.Context, userID int) (*User, error) {
	if err := u.ensureSchema(); err != nil {
		return nil, err
	}

	var user User
	if err := u.db.WithContext(ctx).Where("userid = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user.Sanitized(), nil
}

func (u *UserHandler) ToggleLike(ctx context.Context, userID int, track *music.Objector) (*ToggleLikeResponse, error) {
	if track == nil {
		return nil, errors.New("music track is required")
	}
	if err := u.ensureSchema(); err != nil {
		return nil, err
	}

	var user User
	if err := u.db.WithContext(ctx).Where("userid = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	targetIndex := -1
	for index, item := range user.LikeList {
		if item == nil {
			continue
		}
		if (track.ID != "" && item.ID == track.ID) || (track.FileName != "" && item.FileName == track.FileName) {
			targetIndex = index
			break
		}
	}

	liked := false
	if targetIndex >= 0 {
		user.LikeList = append(user.LikeList[:targetIndex], user.LikeList[targetIndex+1:]...)
	} else {
		trackCopy := *track
		user.LikeList = append(user.LikeList, &trackCopy)
		liked = true
	}

	if err := u.db.WithContext(ctx).Model(&User{}).Where("userid = ?", userID).Update("like_list", user.LikeList).Error; err != nil {
		return nil, err
	}

	user.Password = ""
	return &ToggleLikeResponse{
		Liked:    liked,
		LikeList: user.LikeList,
		User:     user.Sanitized(),
	}, nil
}

func (u *UserHandler) ensureSchema() error {
	return u.db.AutoMigrate(&User{})
}
