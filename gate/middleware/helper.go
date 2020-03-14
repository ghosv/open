package middleware

import (
	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/meta"

	"github.com/ddosakura/sola/v2"

	"github.com/micro/go-micro/v2"
)

// LoadService Middleware Builder
func LoadService(service micro.Service) sola.Middleware {
	client.Default.Init(service)
	return func(h sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			c.Set(meta.CtxService, client.Default)
			return h(c)
		}
	}
}
