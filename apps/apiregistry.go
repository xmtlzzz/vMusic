package apps

import (
	"github.com/emicklei/go-restful/v3"
	musicapi "github.com/xmtlzzz/vMusic/apps/music/api"
	tokenapi "github.com/xmtlzzz/vMusic/apps/token/api"
	userapi "github.com/xmtlzzz/vMusic/apps/user/api"
)

func init() {
	musicWs := musicapi.NewMusicApiHandler().WebService()
	tokenWs := tokenapi.NewTokenApiHandler().WebService()
	userWs := userapi.NewUserApiHandler().WebService()
	restful.Add(musicWs)
	restful.Add(tokenWs)
	restful.Add(userWs)
	registerSwagger()
	registerPlayerUI()
}
