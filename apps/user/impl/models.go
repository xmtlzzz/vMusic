package impl

import (
	music "github.com/xmtlzzz/vMusic/apps/music/impl"
)

type User struct {
	Info UserInfo `json:"info"`

	// 用户喜欢的music切片
	LikeList []*music.Objector `json:"like_list"`

	// 判断用户是否为管理员
	IsAdmin bool `json:"is_admin"`
}

func NewUser() *User {
	var ll []*music.Objector
	return &User{LikeList: ll}
}

type RegistryRequest struct {
	UserInfo User `json:"user_info"`
	// 邀请码
	InvitationCode string `json:"invitation_code"`
}

func NewRegistryRequest() *RegistryRequest {
	return &RegistryRequest{}
}

type UserInfo struct {
	Username string `json:"username"`

	Mail string `json:"mail"`

	Telephone string `json:"telephone"`
}
