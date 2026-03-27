package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	user_impl "github.com/xmtlzzz/vMusic/apps/user/impl"
	"github.com/xmtlzzz/vMusic/utils"
	"gorm.io/gorm"
)

type TokenHandler struct {
	Log zerolog.Logger
	DB  *gorm.DB
}

func NewTokenHandler() *TokenHandler {
	handler := &TokenHandler{
		Log: zerolog.New(os.Stdout),
		DB:  utils.NewDBConnector(),
	}
	if err := handler.DB.AutoMigrate(&Token{}); err != nil {
		handler.Log.Error().Err(err).Msg("failed to migrate token schema")
	}
	return handler
}

var _ TokenInt = (*TokenHandler)(nil)

func (t *TokenHandler) RegistryToken(ctx context.Context, request *RegistryTokenRequest) (*Token, error) {
	if request.RefUserId == 0 {
		return nil, errors.New("ref_user_id is required")
	}

	ns := uuid.NameSpaceURL

	secretKey := uuid.NewSHA1(ns, []byte(fmt.Sprintf("%v.%v", request.Username, request.Password)))
	request.Claims = jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, request.Claims)
	str, err := jwtToken.SignedString(secretKey[:])
	if err != nil {
		log.Printf("Token.SignedString err: %v", err)
		return nil, err
	}
	token := IssueToken(request.RefUserId, str)
	if err := t.DB.AutoMigrate(&Token{}); err != nil {
		return nil, err
	}

	var existed Token
	err = t.DB.WithContext(ctx).Where("ref_user_id = ?", request.RefUserId).First(&existed).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := gorm.G[Token](t.DB).Create(ctx, token); err != nil {
			return nil, err
		}
		return token, nil
	}
	if err != nil {
		return nil, err
	}

	existed.AccessToken = token.AccessToken
	existed.IssueAt = token.IssueAt
	existed.AccessTokenExpireAt = token.AccessTokenExpireAt
	existed.RefreshToken = token.RefreshToken
	existed.RefreshTokenExpireAt = token.RefreshTokenExpireAt
	if err := t.DB.WithContext(ctx).Save(&existed).Error; err != nil {
		return nil, err
	}

	return &existed, nil
}

func (t *TokenHandler) DeleteToken(ctx context.Context, request *DeleteTokenRequest) error {
	user, err := t.QueryUserByAccessToken(ctx, &QueryUserByTokenRequest{AccessToken: request.AccessToken})
	if err != nil {
		return err
	}
	if user.Username != request.Username {
		return errors.New("user does not match token")
	}
	subQuery := t.DB.WithContext(ctx).
		Select("userid").
		Where("username = ?", request.Username).
		Limit(1).Model(&user_impl.User{})
	if err := t.DB.WithContext(ctx).Where("ref_user_id = (?)", subQuery).Delete(&Token{}).Error; err != nil {
		return err
	}
	return nil
}

func IssueToken(refUserID int, token string) *Token {
	expireAt := time.Now().Add(24 * time.Hour)
	return &Token{
		RefUserId:           refUserID,
		AccessToken:         token,
		IssueAt:             time.Now(),
		AccessTokenExpireAt: &expireAt,
	}
}

func (t *TokenHandler) QueryUserByAccessToken(ctx context.Context, request *QueryUserByTokenRequest) (*user_impl.User, error) {
	if request.AccessToken == "" {
		return nil, errors.New("access_token is required")
	}

	var user user_impl.User
	subQuery := t.DB.WithContext(ctx).
		Model(&Token{}).
		Select("ref_user_id").
		Where("access_token = ?", request.AccessToken).
		Limit(1)

	if err := t.DB.WithContext(ctx).Where("userid = (?)", subQuery).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found for token")
		}
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func (t *TokenHandler) CheckValidate(ctx context.Context, token string) (*user_impl.User, error) {
	return t.QueryUserByAccessToken(ctx, &QueryUserByTokenRequest{AccessToken: token})
}
