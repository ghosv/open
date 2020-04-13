package main

import (
	"github.com/ghosv/open/app/mrms/handler"
	pb "github.com/ghosv/open/app/mrms/proto"
	"github.com/ghosv/open/plat/conf"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("srv.mrms"),
		micro.Version("latest"),
	)

	// Connect DB
	mw, md := conf.MongoRepo()
	defer md()

	// Initialize Service
	service.Init(
		micro.WrapHandler(
			mw,
		),
	)

	// Register Handler
	pb.RegisterDeviceHandler(service.Server(), new(handler.Device))
	pb.RegisterRoomHandler(service.Server(), new(handler.Room))
	pb.RegisterMeetingHandler(service.Server(), new(handler.Meeting))

	// Run Service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
