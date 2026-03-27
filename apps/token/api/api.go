package api

import (
	"fmt"
	"os"

	"github.com/emicklei/go-restful/v3"
	"github.com/rs/zerolog"
	token_impl "github.com/xmtlzzz/vMusic/apps/token/impl"
)

type TokenApiHandler struct {
	log zerolog.Logger
}

func NewTokenApiHandler() *TokenApiHandler {
	return &TokenApiHandler{log: zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()}
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

	ws.Route(ws.POST("/remove").To(t.DeleteToken).
		Doc("query user by access token").
		Reads(token_impl.DeleteTokenRequest{}).
		Writes("").
		Returns(200, "OK", "user token delete success").
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
			t.log.Error().Msgf("access_token is required")
			response.WriteErrorString(400, err.Error())
			return
		}
		if err.Error() == "user not found for token" {
			t.log.Error().Msgf("user not found for token")
			response.WriteErrorString(404, err.Error())
			return
		}
		t.log.Error().Msgf("%v", err.Error())
		response.WriteErrorString(500, err.Error())
		return
	}
	t.log.Info().Msgf("query a user: %v", user)
	response.WriteEntity(user)
}

func (t *TokenApiHandler) DeleteToken(request *restful.Request, response *restful.Response) {
	dr := token_impl.NewDeleteTokenRequest()
	if err := request.ReadEntity(dr); err != nil {
		t.log.Error().Msgf("read request value error: %v", err.Error())
		response.WriteErrorString(400, "parse token from request body error")
		return
	}
	if err := token_impl.NewTokenHandler().DeleteToken(request.Request.Context(), dr); err != nil {
		if err.Error() == "user does not match token" {
			t.log.Error().Msgf("user does not match token")
			response.WriteErrorString(400, "user does not match token")
			return
		}
		t.log.Error().Msgf("%v", err.Error())
		response.WriteErrorString(500, err.Error())
		return
	}
	t.log.Info().Msgf("delete user: %v token: %v", dr.Username, dr.AccessToken)
	response.WriteEntity(fmt.Sprintf("delete user: %v token: %v", dr.Username, dr.AccessToken))
	return
}
