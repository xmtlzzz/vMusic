package impl_test

import (
	"context"
	"testing"

	"github.com/xmtlzzz/vMusic/apps/token/impl"
)

func TestTokenHandler_RegistryToken(t *testing.T) {
	tr := impl.NewRegistryTokenRequest()
	tr.Username = "sz"
	tr.Email = "246542922@qq.com"
	tk, err := impl.NewTokenHandler().RegistryToken(context.Background(), tr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)

}
