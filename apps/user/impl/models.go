package impl

import (
	music "github.com/xmtlzzz/vMusic/apps/music/impl"
)

type User struct {
	Userid   int    `json:"userid" gorm:"primary_key;type:int"`
	Username string `json:"username" gorm:"column:username;type:varchar(10);unique"`

	Mail string `json:"mail" gorm:"column:mail;type:varchar(20);unique"`

	Telephone string `json:"telephone" gorm:"column:telephone;type:varchar(11);unique"`

	Password string `json:"password"gorm:"column:password;type:varchar(30)"`

	// 用户喜欢的music切片
	LikeList []*music.Objector `json:"like_list" gorm:"type:json"`

	// 判断用户是否为管理员
	IsAdmin bool `json:"is_admin" gorm:"type:boolean"`
}

func NewUser() *User {
	var ll []*music.Objector
	return &User{LikeList: ll}
}

type CreateUserRequest struct {
	User `json:"user_info"`
	// 邀请码
	InvitationCode string `json:"invitation_code"`
}

func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{}
}

type DeleteUserRequest struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

func NewDeleteUserRequest() *DeleteUserRequest {
	return &DeleteUserRequest{}
}

type QueryUserRequest struct {
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserLoginRequest() *UserLoginRequest {
	return &UserLoginRequest{}
}
