package impl

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/xmtlzzz/vMusic/utils"
	"gorm.io/gorm"
)

type TokenHandler struct {
	Log zerolog.Logger
	DB  *gorm.DB
}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{
		Log: zerolog.New(os.Stdout),
		DB:  utils.NewDBConnector(),
	}
}

var _ TokenInt = (*TokenHandler)(nil)

func (t *TokenHandler) RegistryToken(ctx context.Context, request *RegistryTokenRequest) (*Token, error) {
	ns := uuid.NameSpaceURL

	secretKey := uuid.NewSHA1(ns, []byte(fmt.Sprintf("%v.%v.%v.%v", request.Username, request.Telephone, request.Email)))
	request.Claims = jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, request.Claims)
	str, err := jwtToken.SignedString(secretKey[:])
	if err != nil {
		log.Printf("Token.SignedString err: %v", err)
		return nil, err
	}
	token := IssueToken(str)
	if err := gorm.G[Token](t.DB).Create(ctx, token); err != nil {
		return nil, err
	}
	return token, nil

}

func (t *TokenHandler) DeleteToken(ctx context.Context, request *DeleteTokenRequest) error {
	//TODO implement me
	panic("implement me")
}

func IssueToken(token string) *Token {
	return &Token{
		AccessToken: token,
		IssueAt:     time.Now(),
	}
}
