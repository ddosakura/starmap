package main

import (
	"fmt"

	"github.com/ddosakura/starmap/srv/auth/handler"
	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/auth/subscriber"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name(raw.SrvName),
		micro.Version(raw.SrvVer),
	)

	url := fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		common.ENV.Repo.MySQL.User,
		common.ENV.Repo.MySQL.Pass,
		common.ENV.Repo.MySQL.Host,
		common.ENV.Repo.MySQL.DB,
	)
	// fmt.Println(url)
	repo, e := gorm.Open(
		"mysql",
		url,
	)
	if e != nil {
		log.Fatal(e)
	}
	defer repo.Close()
	repo.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Perm{},
	)

	//log.Log("DB OK!")

	wrapperRepo := common.GormRepo(service, repo)
	wrapperMongo, closeMongo := common.MongoRepo(service)
	defer closeMongo()
	// Initialise service
	service.Init(
		micro.WrapHandler(wrapperRepo),
		micro.WrapHandler(wrapperMongo),
	)

	// Register Handler
	proto.RegisterUserHandler(service.Server(), new(handler.User))
	proto.RegisterRoleHandler(service.Server(), new(handler.Role))
	proto.RegisterPermHandler(service.Server(), new(handler.Perm))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("starmap.srv.auth", service.Server(), new(subscriber.User))
	// Register Function as Subscriber
	micro.RegisterSubscriber("starmap.srv.auth", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
