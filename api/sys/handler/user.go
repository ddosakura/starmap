package handler

import (
	"context"

	"github.com/ddosakura/starmap/api/sys/raw"
	"github.com/micro/go-micro/errors"

	"github.com/ddosakura/starmap/api/rest"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
)

// User Handler
type User struct{}

// Entity API
func (*User) Entity(ctx context.Context, req *api.Request, res *api.Response) error {
	u := new(auth.UserToken)
	u.Auth = new(auth.UserAuth)
	u.User = new(auth.UserInfo)

	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		// API
		// -> username, password
		// <- userinfo
		Action(rest.POST).
		Chain(rest.ParamCheck(rest.PCCS{
			"username": rest.PccMust,
			"password": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(rest.PermCheck([]string{"user:insert"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			userToken, err := s.AuthUserClient.(auth.UserService).Insert(ctx, u.Auth)
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(userToken.User)
		}).
		Done().
		// API
		// -> id
		Action(rest.DELETE).
		Chain(rest.ParamCheck(rest.PCCS{
			"id": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(rest.PermCheck([]string{"user:delete"}, rest.LogicalAND)).
		Chain(rest.RoleLevelCheck(u.Auth.ID)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			_, err := s.AuthUserClient.(auth.UserService).Delete(ctx, u.Auth)
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(nil)
		}).
		Done().
		// API
		// -> id/username
		// <- userinfo
		Action(rest.GET).
		Chain(rest.ParamCheck(rest.PCCS{
			"$auth":    rest.PccLabel(rest.PccMust, rest.LogicalOR),
			"id":       rest.PccLink("$auth"),
			"username": rest.PccLink("$auth"),
		})).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(rest.PermCheck([]string{"user:select"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			userToken, err := s.AuthUserClient.(auth.UserService).Select(ctx, u.Auth)
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(userToken.User)
		}).
		Done().
		// API
		// -> id
		// <- userinfo
		Action(rest.PUT).
		Chain(rest.ParamCheck(rest.PCCS{
			"id": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(rest.ParamAutoLoad(nil, u.User)).
		Chain(rest.PermCheck([]string{"user:update"}, rest.LogicalAND)).
		Chain(rest.RoleLevelCheck(u.Auth.ID)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			if u.Auth.Password == "" {
				u.User.UUID = u.Auth.ID
				u.Auth = nil
			} else {
				u.User = nil
			}
			userToken, err := s.AuthUserClient.(auth.UserService).Change(ctx, u)
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(userToken.User)
		}).
		Done().
		// Finish
		Final()
}

// Role Modify API
func (*User) Role(ctx context.Context, req *api.Request, res *api.Response) error {
	m := new(auth.Modification)
	playload := func(ctx context.Context, s *rest.Flow) error {
		result, err := s.AuthUserClient.(auth.UserService).Role(ctx, m)
		if err != nil {
			return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
		}
		return s.Success(result.Data)
	}

	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		Chain(rest.PermCheck([]string{"user:role"}, rest.LogicalAND)).
		// API - Add role for user
		// -> id, name
		// <- []string
		Action(rest.POST).
		Chain(rest.ParamCheck(rest.PCCS{
			"id":   rest.PccRename(rest.PccMust, "UUID"),
			"name": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, m)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			m.Modify = auth.M_Add
			return nil
		}).
		Chain(playload).
		Done().
		// API - Del role for user
		// -> id, name
		// <- []string
		Action(rest.DELETE).
		Chain(rest.ParamCheck(rest.PCCS{
			"id":   rest.PccRename(rest.PccMust, "UUID"),
			"name": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, m)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			m.Modify = auth.M_Del
			return nil
		}).
		Chain(playload).
		Done().
		// API - All Roles
		// <- []string
		Action(rest.GET).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			m.Modify = auth.M_List
			return nil
		}).
		Chain(playload).
		Done().
		// Finish
		Final()
}
