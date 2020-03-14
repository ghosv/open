package handler

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/router"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/notify/proto"
)

// InjectTempService in /tmp
func InjectTempService(r *router.Router) {
	r.Bind("GET /post_debug", tmpPostDebug)
}

func tmpPostDebug(c sola.Context) error {
	r := c.Request()
	q := r.URL.Query()
	_phone := q.Get("phone")
	_email := q.Get("email")
	_gp := q.Get("gp")

	service := c.Get(meta.CtxService).(*client.MicroClient)
	res, e := service.NotifyVerify.CodeDebug(r.Context(), &pb.VerifyCodeGroup{
		Codes: []*pb.VerifyCode{
			&pb.VerifyCode{
				Type: pb.PostType_CodePhone,
				To:   _phone,
			},
			&pb.VerifyCode{
				Type: pb.PostType_CodeEmail,
				To:   _email,
			},
			&pb.VerifyCode{
				Type: pb.PostType_CodeGP,
				To:   _gp,
			},
		},
	})
	var phoneCode, emailCode, gpCode string
	if e == nil {
		phoneCode = res.Codes[0].Code
		emailCode = res.Codes[1].Code
		gpCode = res.Codes[2].Code
	}

	return c.JSON(http.StatusOK, Response{0, map[string]interface{}{
		"phoneCode": phoneCode,
		"emailCode": emailCode,
		"gpCode":    gpCode,
	}, SUCCESS})
}
