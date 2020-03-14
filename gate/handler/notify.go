package handler

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/utils"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/notify/proto"
)

// NotifyPostAuthCode Handler
func NotifyPostAuthCode(c sola.Context) error {
	r := c.Request()
	q := r.URL.Query()
	_type := q.Get("type") // sms&email&gp(ghost phone)
	_to := q.Get("to")

	var t pb.PostType
	switch _type {
	case "sms":
		t = pb.PostType_CodePhone
	case "email":
		t = pb.PostType_CodeEmail
	case "gp":
		t = pb.PostType_CodeGP
	default:
		return c.JSON(http.StatusOK, Response{-1, nil, FAIL})
	}

	service := c.Get(meta.CtxService).(*client.MicroClient)
	err := service.NotifyPostCode.Publish(r.Context(), &pb.PostCode{
		Type: t,
		To:   _to,
	})
	if err != nil {
		return c.JSON(http.StatusOK, Response{-1, nil,
			utils.MicroErrorDetail(err)})
	}
	return c.JSON(http.StatusOK, Response{0, nil, SUCCESS})
}
