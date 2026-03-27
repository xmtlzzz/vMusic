package api

import (
	"os"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/rs/zerolog"
	music_impl "github.com/xmtlzzz/vMusic/apps/music/impl"
	token_impl "github.com/xmtlzzz/vMusic/apps/token/impl"
	user_impl "github.com/xmtlzzz/vMusic/apps/user/impl"
)

type UserApiHandler struct {
	logger zerolog.Logger
}

type AuthResponse struct {
	Message string            `json:"message"`
	Token   *token_impl.Token `json:"token,omitempty"`
	User    *user_impl.User   `json:"user,omitempty"`
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
		Writes(AuthResponse{}).
		Returns(200, "OK", AuthResponse{}).
		Returns(400, "bad request", nil).
		Returns(500, "internal server error", nil))

	ws.Route(ws.POST("/registry").To(u.UserCreate).
		Doc("create a user").
		Reads(user_impl.CreateUserRequest{}).
		Writes(AuthResponse{}).
		Returns(200, "OK", AuthResponse{}).
		Returns(400, "bad request", nil).
		Returns(500, "internal server error", nil))

	ws.Route(ws.POST("/register").To(u.UserCreate).
		Doc("create a user").
		Reads(user_impl.CreateUserRequest{}).
		Writes(AuthResponse{}).
		Returns(200, "OK", AuthResponse{}).
		Returns(400, "bad request", nil).
		Returns(500, "internal server error", nil))

	ws.Route(ws.GET("/me").To(u.CurrentUser).
		Doc("query current user by access token").
		Writes(user_impl.User{}).
		Returns(200, "OK", user_impl.User{}).
		Returns(401, "unauthorized", nil).
		Returns(500, "internal server error", nil))

	ws.Route(ws.GET("/likes").To(u.ListLikes).
		Doc("list current user's liked music").
		Writes([]music_impl.Objector{}).
		Returns(200, "OK", []music_impl.Objector{}).
		Returns(401, "unauthorized", nil).
		Returns(500, "internal server error", nil))

	ws.Route(ws.POST("/likes/toggle").To(u.ToggleLikeMusic).
		Doc("toggle a music in current user's like list").
		Reads(user_impl.ToggleLikeRequest{}).
		Writes(user_impl.ToggleLikeResponse{}).
		Returns(200, "OK", user_impl.ToggleLikeResponse{}).
		Returns(400, "bad request", nil).
		Returns(401, "unauthorized", nil).
		Returns(500, "internal server error", nil))

	return ws
}

func (u *UserApiHandler) UserCreate(request *restful.Request, response *restful.Response) {
	payload := user_impl.NewCreateUserRequest()
	handler := user_impl.NewUserHandler()
	if err := request.ReadEntity(payload); err != nil {
		u.logger.Warn().Err(err).Msg("info get failed")
		response.WriteErrorString(400, "post参数获取失败")
		return
	}
	if err := handler.CreateUser(request.Request.Context(), payload); err != nil {
		u.logger.Warn().Err(err).Msg("user create failed")
		response.WriteErrorString(400, err.Error())
		return
	}

	user := payload.UserInfo.Sanitized()
	u.logger.Info().Msgf("user: %v is created", user.Username)
	if err := response.WriteEntity(AuthResponse{
		Message: "用户创建完成",
		User:    user,
	}); err != nil {
		response.WriteErrorString(500, err.Error())
	}
}

func (u *UserApiHandler) UserLogin(request *restful.Request, response *restful.Response) {
	userInfo := user_impl.NewUserLoginRequest()
	userHandler := user_impl.NewUserHandler()
	if err := request.ReadEntity(userInfo); err != nil {
		u.logger.Warn().Err(err).Msg("failed to read user info")
		response.WriteErrorString(400, "parse user from Post json error")
		return
	}
	user, err := userHandler.UserIsAvailable(request.Request.Context(), &user_impl.UserLoginRequest{
		Username: userInfo.Username,
		Password: userInfo.Password,
	})
	if err != nil {
		u.logger.Error().Msgf("user unavailable: %v", err)
		response.WriteErrorString(400, err.Error())
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

	if err = response.WriteEntity(AuthResponse{
		Message: "登录成功",
		Token:   token,
		User:    user.Sanitized(),
	}); err != nil {
		response.WriteErrorString(500, err.Error())
	}
}

func (u *UserApiHandler) CurrentUser(request *restful.Request, response *restful.Response) {
	currentUser, err := u.userFromRequest(request)
	if err != nil {
		response.WriteErrorString(401, err.Error())
		return
	}

	if err = response.WriteEntity(currentUser.Sanitized()); err != nil {
		response.WriteErrorString(500, err.Error())
	}
}

func (u *UserApiHandler) ListLikes(request *restful.Request, response *restful.Response) {
	currentUser, err := u.userFromRequest(request)
	if err != nil {
		response.WriteErrorString(401, err.Error())
		return
	}

	if err = response.WriteEntity(currentUser.LikeList); err != nil {
		response.WriteErrorString(500, err.Error())
	}
}

func (u *UserApiHandler) ToggleLikeMusic(request *restful.Request, response *restful.Response) {
	payload := user_impl.NewToggleLikeRequest()
	if err := request.ReadEntity(payload); err != nil {
		response.WriteErrorString(400, "parse like request error")
		return
	}

	if strings.TrimSpace(payload.AccessToken) == "" {
		payload.AccessToken = accessTokenFromRequest(request)
	}
	if strings.TrimSpace(payload.AccessToken) == "" {
		response.WriteErrorString(401, "access token is required")
		return
	}
	if strings.TrimSpace(payload.FileName) == "" {
		response.WriteErrorString(400, "file_name is required")
		return
	}

	currentUser, err := token_impl.NewTokenHandler().QueryUserByAccessToken(request.Request.Context(), &token_impl.QueryUserByTokenRequest{
		AccessToken: payload.AccessToken,
	})
	if err != nil {
		response.WriteErrorString(401, err.Error())
		return
	}

	track, err := music_impl.NewMusicObj().GetTrackByFileName(payload.FileName)
	if err != nil {
		response.WriteErrorString(400, err.Error())
		return
	}

	result, err := user_impl.NewUserHandler().ToggleLike(request.Request.Context(), currentUser.Userid, track)
	if err != nil {
		response.WriteErrorString(500, err.Error())
		return
	}

	if err = response.WriteEntity(result); err != nil {
		response.WriteErrorString(500, err.Error())
	}
}

func (u *UserApiHandler) userFromRequest(request *restful.Request) (*user_impl.User, error) {
	return token_impl.NewTokenHandler().QueryUserByAccessToken(request.Request.Context(), &token_impl.QueryUserByTokenRequest{
		AccessToken: accessTokenFromRequest(request),
	})
}

func accessTokenFromRequest(request *restful.Request) string {
	authorization := strings.TrimSpace(request.Request.Header.Get("Authorization"))
	if authorization != "" {
		if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
			return strings.TrimSpace(authorization[7:])
		}
		return authorization
	}

	return strings.TrimSpace(request.Request.Header.Get("X-Access-Token"))
}
