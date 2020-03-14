package main

import (
	"log"

	"github.com/micro/go-micro/v2"

	"github.com/ghosv/open/gate/conf"
	"github.com/ghosv/open/gate/handler"
	"github.com/ghosv/open/gate/middleware"
	"github.com/ghosv/open/meta"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/native"
	"github.com/ddosakura/sola/v2/middleware/proxy"
	"github.com/ddosakura/sola/v2/middleware/router"
)

func initClient() micro.Service {
	// New Service
	service := micro.NewService(
		micro.Name(meta.ClientGateway),
		micro.Version(meta.Version),
	)

	// Initialise Service
	service.Init()

	return service
}

func main() {
	// TODO: Optional
	cfg := conf.Default

	loadService := middleware.LoadService(initClient())

	sola.Use(proxy.Favicon(cfg.Favicon))
	sola.DefaultApp.LoadConfig()

	{
		r := router.New(nil)
		r.Bind("GET /", handler.HomePage)
		r.Use(loadService)

		// notify
		// TODO: 发送 code 注意限频率
		r.Bind("GET /notify/post/code", handler.NotifyPostAuthCode)

		{
			sub := r.Sub(&router.Option{Pattern: "/tmp"})
			handler.InjectTempService(sub)
		}

		// oauth2
		r.Bind("GET /authorize", handler.Authorize) // TODO: 兼容已登录状态下免填token
		r.Bind("GET /authorize_confirm", handler.AuthorizeConfirm)
		r.Bind("POST /access_token", handler.AccessToken)

		r.Bind("POST /register", handler.Register)
		r.Bind("POST /login", handler.Login)
		// r.Bind("/logout", h)

		r.Use(middleware.Auth)
		// Check Login Status
		r.Bind("/status", handler.Status)
		{
			sub := r.Sub(&router.Option{Pattern: "/tmp2"})
			handler.InjectTempService2(sub)
		}

		r.Bind("POST /graphql", handler.DefaultGraphQL)

		{
			sub := r.Sub(&router.Option{Pattern: "/api"})
			// TODO
			sub.Bind("/test", nil)
		}

		sola.Use(r.Routes())
	}

	sola.Use(native.Static("static", ""))

	log.Println(conf.ReadyServer)
	sola.ListenKeep(cfg.Addr)
}
