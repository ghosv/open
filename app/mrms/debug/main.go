package main

import (
	"context"

	pb "github.com/ghosv/open/app/mrms/proto"
	"github.com/kr/pretty"
	"github.com/micro/go-micro/v2"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("debug.mrms"),
		micro.Version("latest"),
	)

	// Initialize Service
	service.Init(
	// TODO: 预处理？
	)

	// Register Handler
	c := pb.NewDeviceService("srv.mrms", service.Client())
	ctx := context.Background()

	// pretty.Println(create(ctx, c))
	pretty.Println(find(ctx, c))
	pretty.Println(search(ctx, c))
	// pretty.Println(delete(ctx, c))
	pretty.Println(update(ctx, c))
}

func create(ctx context.Context, c pb.DeviceService) (interface{}, error) {
	return c.Create(ctx, &pb.DeviceModify{
		Info: &pb.DeviceInfo{
			Name: "abc3",
			Type: "yyz",
		},
	})
}

func find(ctx context.Context, c pb.DeviceService) (interface{}, error) {
	return c.BatchFind(ctx, &pb.BatchID{
		UUID: []string{"5e93404a1d4f591abb73c616"}})
}

func search(ctx context.Context, c pb.DeviceService) (interface{}, error) {
	return c.Search(ctx, &pb.SearchForm{
		Page: 1,
		Size: 5,
		Word: "b",
	})
}

func delete(ctx context.Context, c pb.DeviceService) (interface{}, error) {
	return c.Delete(ctx, &pb.DeviceModify{Info: &pb.DeviceInfo{
		ID: "5e93404a1d4f591abb73c616"}})
}

func update(ctx context.Context, c pb.DeviceService) (interface{}, error) {
	return c.Update(ctx, &pb.DeviceModify{Info: &pb.DeviceInfo{
		ID:   "5e9340671d4f591abb73c617",
		Name: "dc",
	}})
}
