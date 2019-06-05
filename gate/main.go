package main

import (
	"net/http"

	"github.com/ddosakura/starmap/gate/middleware/client"
	"github.com/ddosakura/starmap/gate/router"
	"github.com/labstack/echo"

	_ "github.com/ddosakura/starmap/gate/check/echo"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("starmap.gate"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	e := echo.New()
	e.Use(client.AuthWrapper(service))

	router.Bind(e)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	go func() {
		e.Logger.Fatal(e.Start(ENV.Gateway.Addr))
	}()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
