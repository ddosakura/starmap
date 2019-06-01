package client

import (
	"context"

	"github.com/ddosakura/starmap/api/common"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	authRaw "github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
)

// Wrapper returns a wrapper for the Client
func Wrapper(service micro.Service) server.HandlerWrapper {
	userClient := auth.NewUserService(authRaw.SrvName, service.Client())
	roleClient := auth.NewRoleService(authRaw.SrvName, service.Client())
	permClient := auth.NewPermService(authRaw.SrvName, service.Client())
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, common.SrvKey{Srv: common.SrvAuthUser}, userClient)
			ctx = context.WithValue(ctx, common.SrvKey{Srv: common.SrvAuthRole}, roleClient)
			ctx = context.WithValue(ctx, common.SrvKey{Srv: common.SrvAuthPerm}, permClient)
			return fn(ctx, req, rsp)
		}
	}
}
