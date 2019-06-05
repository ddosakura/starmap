package client

import (
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	raw "github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/labstack/echo"
	"github.com/micro/go-micro"
)

const (
	srvAuthUser = "srv.auth.user"
	srvAuthRole = "srv.auth.role"
	srvAuthPerm = "srv.auth.perm"
)

// AuthWrapper for Client
func AuthWrapper(service micro.Service) echo.MiddlewareFunc {
	userClient := auth.NewUserService(raw.SrvName, service.Client())
	roleClient := auth.NewRoleService(raw.SrvName, service.Client())
	permClient := auth.NewPermService(raw.SrvName, service.Client())
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(srvAuthUser, userClient)
			ctx.Set(srvAuthRole, roleClient)
			ctx.Set(srvAuthPerm, permClient)
			return next(ctx)
		}
	}
}

// AuthUser from ctx
func AuthUser(ctx echo.Context) auth.UserService {
	return ctx.Get(srvAuthUser).(auth.UserService)
}

// AuthRole from ctx
func AuthRole(ctx echo.Context) auth.RoleService {
	return ctx.Get(srvAuthRole).(auth.RoleService)
}

// AuthPerm from ctx
func AuthPerm(ctx echo.Context) auth.PermService {
	return ctx.Get(srvAuthPerm).(auth.PermService)
}
