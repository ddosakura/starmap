package main

import (
	"github.com/ddosakura/starmap/srv/auth/handler"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/auth/subscriber"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name(raw.SrvName),
		micro.Version(raw.SrvVer),
	)

	wrapperRepo, closeRepo := common.MongoRepo(service)
	defer closeRepo()
	// Initialise service
	service.Init(
		micro.WrapHandler(wrapperRepo),
	)

	// Register Handler
	proto.RegisterUserHandler(service.Server(), new(handler.User))
	proto.RegisterRoleHandler(service.Server(), new(handler.Role))
	proto.RegisterPermissionHandler(service.Server(), new(handler.Permission))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("starmap.srv.auth", service.Server(), new(subscriber.User))
	// Register Function as Subscriber
	micro.RegisterSubscriber("starmap.srv.auth", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
