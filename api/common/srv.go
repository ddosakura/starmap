package common

import (
	"context"

	"github.com/ddosakura/starmap/api/rest"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
)

type srv int

// Srv for key
const (
	SrvAuthUser srv = iota
	SrvAuthRole
	SrvAuthPerm
)

// SrvKey for Ctx
type SrvKey struct {
	Srv srv
}

// AuthServiceUserFromContext retrieves the client from the Context
func AuthServiceUserFromContext(ctx context.Context) (rest.AuthService, bool) {
	c, ok := ctx.Value(SrvKey{SrvAuthUser}).(auth.UserService)
	return c, ok
}

// AuthUserFromContext retrieves the client from the Context
func AuthUserFromContext(ctx context.Context) (auth.UserService, bool) {
	c, ok := ctx.Value(SrvKey{SrvAuthUser}).(auth.UserService)
	return c, ok
}

// AuthRoleFromContext retrieves the client from the Context
func AuthRoleFromContext(ctx context.Context) (auth.RoleService, bool) {
	c, ok := ctx.Value(SrvKey{SrvAuthRole}).(auth.RoleService)
	return c, ok
}

// AuthPermFromContext retrieves the client from the Context
func AuthPermFromContext(ctx context.Context) (auth.PermService, bool) {
	c, ok := ctx.Value(SrvKey{SrvAuthPerm}).(auth.PermService)
	return c, ok
}
