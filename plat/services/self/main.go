package main

import (
	"github.com/ghosv/open/meta"
	notifyPb "github.com/ghosv/open/plat/services/notify/proto"
	"github.com/ghosv/open/plat/services/self/handler"
	"github.com/ghosv/open/plat/services/self/initialize"
	pb "github.com/ghosv/open/plat/services/self/proto"
	"github.com/ghosv/open/plat/utils"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name(initialize.SrvName),
		micro.Version(initialize.SrvVer),
	)

	// Connect DB
	repo := initialize.Repo()
	defer repo.Close()

	// Initialize Service
	service.Init(
		micro.WrapHandler(
			// TODO: remove debug
			utils.MicroHandlerContextWithValue(meta.KeyRepo, repo.Debug()),

			utils.MicroHandlerContextWithValue(meta.KeyClientNotifyVerify,
				notifyPb.NewVerifyService(meta.SrvNotify, service.Client())),
		),
	)

	// Register Handler
	pb.RegisterUserHandler(service.Server(), new(handler.User))
	pb.RegisterAppHandler(service.Server(), new(handler.App))
	pb.RegisterOrgHandler(service.Server(), new(handler.Org))
	pb.RegisterGroupHandler(service.Server(), new(handler.Group))

	// Run Service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
