package main

import (
	"github.com/micro/go-log"

	"github.com/ddosakura/starmap/api/auth/client"
	"github.com/ddosakura/starmap/api/auth/handler"
	proto "github.com/ddosakura/starmap/api/auth/proto"
	"github.com/ddosakura/starmap/api/auth/raw"
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

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
