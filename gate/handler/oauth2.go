package handler

import (
	"net/http"
	"strings"

	"github.com/ddosakura/sola/v2"
	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/utils"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/core/proto"
)

// Authorize GET /authorize
func Authorize(c sola.Context) error {
	r := c.Request()
	query := r.URL.Query()
	appID := query.Get("id")
	token := query.Get("token")
	redirect := query.Get("redirect")

	if appID == "" {
		// TODO: 美化
		return c.HTML(http.StatusOK, `<div>Unknow App</div>`)
	}
	if token == "" {
		// TODO: 美化
		return c.HTML(http.StatusOK, `<div>Logout</div>`)
	}

	service := c.Get(meta.CtxService).(*client.MicroClient)

	res, err := service.CoreAuth.Authorize(r.Context(), &pb.AuthRequest{
		Type:  pb.AuthType_Code,
		Token: token,
		AppID: appID,
	})
	if err != nil {
		// TODO: 美化
		return c.HTML(http.StatusOK, `<div>Fail</div>`+err.Error())
	}

	scopes := ""
	txt := ""
	for i, v := range res.Info.AccessList {
		if i > 0 {
			scopes += "|"
		}
		scopes += v.Name
		txt += v.Detail + "; "
	}

	redirect = "/authorize_confirm?code=" + res.Code +
		"&redirect=" + redirect +
		"&scopes=" + scopes

	// TODO: 授权页美化 & 在页面上选择 scopes 参数
	return c.HTML(http.StatusOK, `
		<div>申请权限：`+txt+`</div>
		<div>Confirm?</div>
		<a href="`+redirect+`">yes</a>
	`)
}

// AuthorizeConfirm GET /authorize_confirm
func AuthorizeConfirm(c sola.Context) error {
	r := c.Request()
	query := r.URL.Query()
	code := query.Get("code")
	_scopes := query.Get("scopes")
	redirect := query.Get("redirect")
	var scopes []*pb.AccessScope
	if _scopes == "" {
		scopes = []*pb.AccessScope{}
	} else {
		scopeArray := strings.Split(_scopes, "|")
		scopes = make([]*pb.AccessScope, 0, len(scopeArray))
		for _, v := range scopeArray {
			scopes = append(scopes, &pb.AccessScope{Name: v})
		}
	}

	service := c.Get(meta.CtxService).(*client.MicroClient)
	_, err := service.CoreAuth.AuthorizeConfirm(r.Context(), &pb.AuthConfirm{
		Code:   code,
		Scopes: scopes,
	})
	if err != nil {
		// TODO: 美化
		return c.HTML(http.StatusOK, `<div>Fail</div>`+err.Error())
	}

	redirect = redirect + "?code=" + code

	w := c.Response()
	w.Header().Add("Location", redirect)
	w.WriteHeader(http.StatusMovedPermanently)
	return nil
}

// AccessToken POST /access_token
func AccessToken(c sola.Context) error {
	r := c.Request()
	r.ParseForm()
	appID := r.PostFormValue("app_id")
	appKey := r.PostFormValue("app_key")
	appSecret := r.PostFormValue("app_secret")
	code := r.PostFormValue("code")

	service := c.Get(meta.CtxService).(*client.MicroClient)
	res, err := service.CoreAuth.AuthToken(r.Context(), &pb.AuthApp{
		AppID:     appID,
		AppKey:    appKey,
		AppSecret: appSecret,
		Code:      code,
	})
	if err != nil {
		return c.JSON(http.StatusOK, Response{-1, nil,
			utils.MicroErrorDetail(err)})
	}

	return c.JSON(http.StatusOK, Response{0, map[string]interface{}{
		"token": res.Str,
	}, SUCCESS})
}
