package middleware

import (
	"net/http"
	"strings"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/handler"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/core/proto"

	"github.com/ddosakura/sola/v2"
)

// Auth Middleware
func Auth(next sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		r := c.Request()
		path := r.URL.Path

		if path != "/status" &&
			path != "/graphql" &&
			!strings.HasPrefix(path, "/api") &&
			!strings.HasPrefix(path, "/tmp2") {
			return next(c)
		}

		token := r.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1) // TODO: 前端注意对接
		service := c.Get(meta.CtxService).(*client.MicroClient)
		res, err := service.CoreUser.Check(r.Context(), &pb.Token{
			Str: token,
		})
		if err != nil {
			if path == "/graphql" {
				return c.JSON(http.StatusUnauthorized, nil)
			}
			return c.JSON(http.StatusOK, handler.Response{
				Code: -1,
				Msg:  handler.LOGOUT,
			})
		}

		// TODO: 前端注意对接刷新 Token，为空字符串时代表未刷新
		c.Response().Header().Set("Set-Token", res.Str)
		c.Set(meta.CtxTokenPayload, res.Payload)
		return next(c)
	}
}
