package handler

import (
	"net/http"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/utils"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/core/proto"

	"github.com/ddosakura/sola/v2"
)

type loginForm struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Pass  string `json:"pass"`

	PhoneCode string `json:"phoneCode"`
	EmailCode string `json:"emailCode"`
}

// Register POST /register
func Register(c sola.Context) error {
	d := loginForm{}
	c.GetJSON(&d)

	service := c.Get(meta.CtxService).(*client.MicroClient)

	r := c.Request()
	res, err := service.CoreUser.Register(r.Context(), &pb.Identity{
		Name: d.Name,
		Pass: d.Pass,

		Phone:     d.Phone,
		Email:     d.Email,
		PhoneCode: d.PhoneCode,
		EmailCode: d.EmailCode,

		// TODO: more info
	})
	if err != nil {
		return c.JSON(http.StatusOK, Response{-1, nil,
			utils.MicroErrorDetail(err)})
	}

	return c.JSON(http.StatusOK, Response{0, map[string]string{
		"token": res.Str,
	}, SUCCESS})
}

// Login POST /login
func Login(c sola.Context) error {
	d := loginForm{}
	c.GetJSON(&d)

	service := c.Get(meta.CtxService).(*client.MicroClient)

	r := c.Request()
	res, err := service.CoreUser.Login(r.Context(), &pb.Credential{
		NamePhoneEmail: d.Name,
		Pass:           d.Pass,
	})
	if err != nil {
		return c.JSON(http.StatusOK, Response{-1, nil,
			utils.MicroErrorDetail(err)})
	}

	return c.JSON(http.StatusOK, Response{0, map[string]string{
		"token": res.Str,
	}, SUCCESS})
}

// TODO: forget password & reset password
// resetUserPass(UUID: String!, type: AuthType!, to: String!, code: String!, pass: String!): Boolean!
