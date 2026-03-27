package impl_test

import (
	"context"
	"testing"

	"github.com/xmtlzzz/vMusic/apps/user/impl"
)

func TestCreateUser(t *testing.T) {
	user := impl.NewCreateUserRequest()
	user.UserInfo.Username = "sz"
	user.UserInfo.Telephone = "13136325596"
	user.UserInfo.Mail = "2465429244@qq.com"
	user.UserInfo.Password = "123456"
	if err := impl.NewUserHandler().CreateUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}
	t.Logf("create user success")
}

func TestDeleteUser(t *testing.T) {
	user := impl.NewDeleteUserRequest()
	user.Username = "test"
	nu, err := impl.NewUserHandler().DeleteUser(context.Background(), &impl.DeleteUserRequest{Username: user.Username})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("delete user success num: %d", nu)
}

func TestUserIsAvailable(t *testing.T) {
	request := impl.NewUserLoginRequest()
	request.Username = "sz"
	request.Password = "123456"
	user, err := impl.NewUserHandler().UserIsAvailable(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user is available: %v", user)
}
