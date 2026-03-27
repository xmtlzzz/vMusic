package impl

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	// 作为数据表的主键
	Id int `json:"id" gorm:"id;primaryKey"`
	// 用户id
	RefUserId int `json:"ref_user_id" gorm:"column:ref_user_id;unique;index"`
	// 颁发的Access token，token要唯一
	AccessToken string `json:"access_token" gorm:"column:access_token;index"`
	// 颁发时间，不能为nil
	IssueAt time.Time `json:"issue_at" gorm:"column:issue_at"`
	// 指针，因为token第一次颁发没有过期可以为nil
	AccessTokenExpireAt *time.Time `json:"access_token_expire_at" gorm:"column:access_token_expire_at"`
	// 刷新token，解决user session问题模拟长连接，tokn要唯一
	RefreshToken string `json:"refresh_token" gorm:"column:refresh_token;unique;index"`
	// 同上因为第一次颁发可以不过期
	RefreshTokenExpireAt *time.Time `json:"refresh_token_expire_at" gorm:"column:refresh_token_expire_at"`
	// 不作为字段插入数据库，但是希望可以进行关联查询：通过RefUserId找到UserName并返回
	RefUserName string `json:"ref_user_name" gorm:"-"`
}

func NewToken() *Token {
	return &Token{}
}

type RegistryTokenRequest struct {
	Username  string               `json:"username"`
	Password  string               `json:"password"`
	RefUserId int                  `json:"ref_user_id"`
	Claims    jwt.RegisteredClaims `json:"claims"`
}

func NewRegistryTokenRequest() *RegistryTokenRequest {
	return &RegistryTokenRequest{}
}

type QueryUserByTokenRequest struct {
	AccessToken string `json:"access_token"`
}

func NewQueryUserByTokenRequest() *QueryUserByTokenRequest {
	return &QueryUserByTokenRequest{}
}

type DeleteTokenRequest struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

func NewDeleteTokenRequest() *DeleteTokenRequest {
	return &DeleteTokenRequest{}
}
