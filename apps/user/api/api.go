package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/rs/zerolog"
	token_impl "github.com/xmtlzzz/vMusic/apps/token/impl"
	user_impl "github.com/xmtlzzz/vMusic/apps/user/impl"
)

type UserApiHandler struct {
	logger zerolog.Logger
}

func (u *UserApiHandler) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/user").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/login").To(u.UserLogin).
		// 设置接收类型可以是multipart/form-data
		Consumes("multipart/form-data").
		Doc("upload music files to local path").
		Returns(200, "OK", map[string]interface{}{}).
		Returns(400, "bad request", nil).
		Returns(500, "internal server error", nil))

	return ws
}

func (u *UserApiHandler) UserLogin(request *restful.Request, response *restful.Response) {
	userInfo := user_impl.NewCreateUserRequest()
	uH := user_impl.NewUserHandler()
	if err := request.ReadEntity(userInfo); err != nil {
		u.logger.Warn().Err(err).Msg("failed to read user info")
		response.WriteErrorString(500, "parse user from Post json error")
		return
	}
	_, err := uH.UserIsAvailable(request.Request.Context(), &user_impl.UserLoginRequest{Username: userInfo.Username, Password: userInfo.Password})
	if err != nil {
		if err.Error() == "user not found" {
			if err := uH.CreateUser(request.Request.Context(), userInfo); err != nil {
				u.logger.Warn().Err(err).Msgf("user create fail")
				response.WriteErrorString(500, "user create failed,login error")
				return
			}
		} else {
			u.logger.Warn().Err(err).Msgf("user unavailable: %v", err)
			response.WriteErrorString(500, "checker user is available error")
			return
		}
	}
	token, err := token_impl.NewTokenHandler().RegistryToken(request.Request.Context(), &token_impl.RegistryTokenRequest{
		Username:  userInfo.Username,
		Email:     userInfo.Mail,
		Telephone: userInfo.Telephone,
	})
	if err != nil {
		u.logger.Warn().Err(err).Msgf("issue token api error: %v", err)
		response.WriteErrorString(500, "user login failed: issue token error")
		return
	}
	response.WriteEntity(token)
}
