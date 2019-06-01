package main

import (
	"github.com/micro/go-log"

	"github.com/ddosakura/starmap/api/sys/client"
	"github.com/ddosakura/starmap/api/sys/handler"
	proto "github.com/ddosakura/starmap/api/sys/proto"
	"github.com/ddosakura/starmap/api/sys/raw"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name(raw.SrvName),
		micro.Version(raw.SrvVer),
	)

	// Initialise service
	service.Init(
		// create wrap for the srv client
		micro.WrapHandler(client.Wrapper(service)),
	)

	// Register Handler
	proto.RegisterUserHandler(service.Server(), new(handler.User))
	proto.RegisterRoleHandler(service.Server(), new(handler.Role))
	proto.RegisterPermHandler(service.Server(), new(handler.Perm))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
