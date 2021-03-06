package handler

import (
	"context"

	"github.com/ddosakura/starmap/api/auth/raw"
	"github.com/ddosakura/starmap/api/common"
	"github.com/ddosakura/starmap/api/rest"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// Roles API
func (*User) Roles(ctx context.Context, req *api.Request, res *api.Response) error {
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		// API - 自身角色查询
		// <- []string
		Action(rest.POST | rest.GET).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			res, err := s.AuthUserClient.Roles(s.Ctx, &auth.Identity{
				UUID: s.Token.User.UUID,
			})
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(res.Data)
		}).
		Done().
		// Finish
		Final()
}

// Perms API
func (*User) Perms(ctx context.Context, req *api.Request, res *api.Response) error {
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		// API - 自身权限查询
		// <- []string
		Action(rest.POST | rest.GET).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			res, err := s.AuthUserClient.Perms(s.Ctx, &auth.Identity{
				UUID: s.Token.User.UUID,
			})
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(res.Data)
		}).
		Done().
		// Finish
		Final()
}

// Perms API
func (*Role) Perms(ctx context.Context, req *api.Request, res *api.Response) error {
	authRoleClient, ok := common.AuthRoleFromContext(ctx)
	if !ok {
		return errors.InternalServerError(raw.SrvName, "auth client not found")
	}

	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		// 系统管理员专用
		Chain(rest.RoleCheck([]string{"sys"}, rest.LogicalAND)).
		// API - 角色权限查询
		// <- []string
		Action(rest.POST | rest.GET).
		Chain(rest.ParamCheck(rest.PCCS{
			"role": &rest.PCC{Must: true},
		})).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			res, err := authRoleClient.Perms(s.Ctx, &auth.Identity{
				Name: s.Params["role"].Values[0],
			})
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(res.Data)
		}).
		Done().
		// Finish
		Final()
}
