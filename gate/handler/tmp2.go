package handler

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/router"
	"github.com/ghosv/open/meta"
)

// InjectTempService2 in /tmp2 <need login>
func InjectTempService2(r *router.Router) {
	r.Bind("GET /token_debug", tmpTokenDebug)
}

func tmpTokenDebug(c sola.Context) error {
	return c.JSON(http.StatusOK, c.Get(meta.CtxTokenPayload))
}
