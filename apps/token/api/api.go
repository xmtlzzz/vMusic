package api

import (
	"github.com/emicklei/go-restful/v3"
	token_impl "github.com/xmtlzzz/vMusic/apps/token/impl"
)

type TokenApiHandler struct {
}

func NewTokenApiHandler() *TokenApiHandler {
	return &TokenApiHandler{}
}

func (t *TokenApiHandler) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/token").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Token APIs")

	ws.Route(ws.POST("/user").To(t.QueryUserByAccessToken).
		Doc("query user by access token").
		Reads(token_impl.QueryUserByTokenRequest{}).
		Writes(map[string]any{}).
		Returns(200, "OK", map[string]any{}).
		Returns(400, "bad request", nil).
		Returns(404, "not found", nil).
		Returns(500, "internal server error", nil))

	return ws
}

func (t *TokenApiHandler) QueryUserByAccessToken(request *restful.Request, response *restful.Response) {
	query := token_impl.NewQueryUserByTokenRequest()
	if err := request.ReadEntity(query); err != nil {
		response.WriteErrorString(400, "parse token from request body error")
		return
	}
	user, err := token_impl.NewTokenHandler().QueryUserByAccessToken(request.Request.Context(), query)
	if err != nil {
		if err.Error() == "access_token is required" {
			response.WriteErrorString(400, err.Error())
			return
		}
		if err.Error() == "user not found for token" {
			response.WriteErrorString(404, err.Error())
			return
		}
		response.WriteErrorString(500, err.Error())
		return
	}

	response.WriteEntity(user)
}
