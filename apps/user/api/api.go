package api

import (
	"os"

	"github.com/emicklei/go-restful/v3"
	"github.com/rs/zerolog"
	token_impl "github.com/xmtlzzz/vMusic/apps/token/impl"
	user_impl "github.com/xmtlzzz/vMusic/apps/user/impl"
)

type UserApiHandler struct {
	logger zerolog.Logger
}

func NewUserApiHandler() *UserApiHandler {
	return &UserApiHandler{
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
}

func (u *UserApiHandler) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/user").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("User APIs")

	ws.Route(ws.POST("/login").To(u.UserLogin).
		Doc("login a user and issue token").
		Reads(user_impl.UserLoginRequest{}).
		Writes(token_impl.Token{}).
		Returns(200, "OK", token_impl.Token{}).
		Returns(400, "bad request", nil).
		Returns(500, "internal server error", nil))

	ws.Route(ws.POST("/registry").To(u.UserCreate).
		Doc("create a user").
		Reads(user_impl.CreateUserRequest{}).
		Writes("").
		Returns(200, "OK", "user create success").
		Returns(400, "bad request", "user create failed").
		Returns(500, "internal server error", "server internal error"))

	return ws
}

func (u *UserApiHandler) UserCreate(request *restful.Request, response *restful.Response) {
	ur := user_impl.NewCreateUserRequest()
	uh := user_impl.NewUserHandler()
	if err := request.ReadEntity(ur); err != nil {
		u.logger.Warn().Err(err).Msg("info get failed")
		response.WriteErrorString(500, "post参数获取失败")
		return
	}
	if err := uh.CreateUser(request.Request.Context(), ur); err != nil {
		u.logger.Warn().Err(err).Msg("user create failed")
		response.WriteErrorString(500, "用户已经存在")
		return
	}
	u.logger.Info().Msgf("user: %v is created", ur.Username)
	response.WriteErrorString(500, "用户创建完成")
	return
}

func (u *UserApiHandler) UserLogin(request *restful.Request, response *restful.Response) {
	userInfo := user_impl.NewUserLoginRequest()
	uH := user_impl.NewUserHandler()
	if err := request.ReadEntity(userInfo); err != nil {
		u.logger.Warn().Err(err).Msg("failed to read user info")
		response.WriteErrorString(500, "parse user from Post json error")
		return
	}
	user, err := uH.UserIsAvailable(request.Request.Context(), &user_impl.UserLoginRequest{Username: userInfo.Username, Password: userInfo.Password})
	if err != nil {
		u.logger.Warn().Err(err).Msgf("user unavailable: %v", err)
		response.WriteErrorString(500, "checker user is available error")
		return
	}
	token, err := token_impl.NewTokenHandler().RegistryToken(request.Request.Context(), &token_impl.RegistryTokenRequest{
		Username:  userInfo.Username,
		Password:  userInfo.Password,
		RefUserId: user.Userid,
	})
	if err != nil {
		u.logger.Warn().Err(err).Msgf("issue token api error: %v", err)
		response.WriteErrorString(500, "user login failed: issue token error")
		return
	}
	response.WriteEntity(token)
}
