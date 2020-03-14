package main

import (
	"github.com/ghosv/open/meta"
	"github.com/ghosv/open/plat/services/notify/handler"
	"github.com/ghosv/open/plat/services/notify/initialize"
	pb "github.com/ghosv/open/plat/services/notify/proto"
	"github.com/ghosv/open/plat/services/notify/subscriber"
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
	repoRedis := initialize.Redis()
	defer repoRedis.Close()

	// Initialize Service
	service.Init(
		micro.WrapHandler(
			utils.MicroHandlerContextWithValue(meta.KeyRepoRedis, repoRedis),
		),
		micro.WrapSubscriber(
			utils.MicroSubscriberContextWithValue(meta.KeyRepoRedis, repoRedis),
		),
	)

	// Register Handler
	pb.RegisterVerifyHandler(service.Server(), new(handler.Verify))

	// Register Struct as Subscriber
	micro.RegisterSubscriber(meta.TopicNotifyPostCode, service.Server(), new(subscriber.PostCode))

	// Run Service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
