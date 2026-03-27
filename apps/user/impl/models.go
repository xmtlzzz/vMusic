package impl

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	music "github.com/xmtlzzz/vMusic/apps/music/impl"
)

type LikeList []*music.Objector

func (l LikeList) Value() (driver.Value, error) {
	if l == nil {
		return "[]", nil
	}

	payload, err := json.Marshal([]*music.Objector(l))
	if err != nil {
		return nil, err
	}
	return string(payload), nil
}

func (l *LikeList) Scan(value any) error {
	if l == nil {
		return nil
	}

	switch data := value.(type) {
	case nil:
		*l = LikeList{}
		return nil
	case []byte:
		return l.unmarshal(data)
	case string:
		return l.unmarshal([]byte(data))
	default:
		return nil
	}
}

func (l *LikeList) unmarshal(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*l = LikeList{}
		return nil
	}

	if bytes.Equal(trimmed, []byte("{}")) {
		*l = LikeList{}
		return nil
	}

	if trimmed[0] == '[' {
		var list []*music.Objector
		if err := json.Unmarshal(trimmed, &list); err != nil {
			return err
		}
		*l = LikeList(list)
		return nil
	}

	var single music.Objector
	if err := json.Unmarshal(trimmed, &single); err != nil {
		return err
	}

	if single == (music.Objector{}) {
		*l = LikeList{}
		return nil
	}

	*l = LikeList{&single}
	return nil
}

type User struct {
	Userid   int    `json:"userid" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"column:username;type:varchar(32);uniqueIndex;not null"`

	Mail string `json:"mail" gorm:"column:mail;type:varchar(64);uniqueIndex"`

	Telephone string `json:"telephone" gorm:"column:telephone;type:varchar(20);uniqueIndex"`

	Password string `json:"password,omitempty" gorm:"column:password;type:varchar(128);not null"`

	// 用户喜欢的music切片
	LikeList LikeList `json:"like_list" gorm:"type:jsonb"`

	// 判断用户是否为管理员
	IsAdmin bool `json:"is_admin" gorm:"type:boolean"`
}

func NewUser() *User {
	var ll LikeList
	return &User{LikeList: ll}
}

type CreateUserRequest struct {
	UserInfo User `json:"user_info"`
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

type ToggleLikeRequest struct {
	AccessToken string `json:"access_token"`
	FileName    string `json:"file_name"`
}

func NewToggleLikeRequest() *ToggleLikeRequest {
	return &ToggleLikeRequest{}
}

type ToggleLikeResponse struct {
	Liked    bool     `json:"liked"`
	LikeList LikeList `json:"like_list"`
	User     *User    `json:"user"`
}

func (u *User) Sanitized() *User {
	if u == nil {
		return nil
	}

	clone := *u
	clone.Password = ""
	if clone.LikeList == nil {
		clone.LikeList = LikeList{}
	}
	return &clone
}
