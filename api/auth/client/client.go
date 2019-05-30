package client

import (
	"context"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	authRaw "github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
)

type srv int

const (
	srvAuthUser srv = iota
)

type key struct {
	srv srv
}

// Wrapper returns a wrapper for the Client
func Wrapper(service micro.Service) server.HandlerWrapper {
	userClient := auth.NewUserService(authRaw.SrvName, service.Client())
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, key{srvAuthUser}, userClient)
			return fn(ctx, req, rsp)
		}
	}
}

// AuthUserFromContext retrieves the client from the Context
func AuthUserFromContext(ctx context.Context) (auth.UserService, bool) {
	c, ok := ctx.Value(key{srvAuthUser}).(auth.UserService)
	return c, ok
}
