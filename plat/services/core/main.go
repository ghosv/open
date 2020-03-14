package main

import (
	"github.com/ghosv/open/meta"
	"github.com/ghosv/open/plat/services/core/handler"
	"github.com/ghosv/open/plat/services/core/initialize"
	pb "github.com/ghosv/open/plat/services/core/proto"
	notifyPb "github.com/ghosv/open/plat/services/notify/proto"
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
	repoRedis := initialize.Redis()
	defer repoRedis.Close()

	// Initialize Service
	service.Init(
		micro.WrapHandler(
			// TODO: remove debug
			utils.MicroHandlerContextWithValue(meta.KeyRepo, repo.Debug()),
			utils.MicroHandlerContextWithValue(meta.KeyRepoRedis, repoRedis),

			utils.MicroHandlerContextWithValue(meta.KeyClientNotifyVerify,
				notifyPb.NewVerifyService(meta.SrvNotify, service.Client())),
		),
	)

	// Register Handler
	pb.RegisterUserHandler(service.Server(), new(handler.User))
	pb.RegisterAuthHandler(service.Server(), new(handler.Auth))
	pb.RegisterAppHandler(service.Server(), new(handler.App))

	// Run Service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
