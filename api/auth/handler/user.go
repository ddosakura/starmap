package handler

import (
	"context"

	"github.com/ddosakura/starmap/api/auth/raw"
	"github.com/ddosakura/starmap/api/rest"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// User Handler
type User struct {
}

// Login API
func (*User) Login(ctx context.Context, req *api.Request, res *api.Response) error {
	a := new(auth.UserAuth)
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		// API - 登录
		// -> username, password
		// <- token, user
		// !! set-token
		Action(rest.POST | rest.GET).
		Chain(rest.ParamCheck(map[string]*rest.PCC{
			"username": rest.PccMust,
			"password": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, a)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			userToken, err := s.AuthUserClient.(auth.UserService).Login(ctx, a)
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			s.FreshJWT(userToken.Token)
			return s.Success(userToken)
		}).
		Done().
		// Finish
		Final()
}

// Register API
func (*User) Register(ctx context.Context, req *api.Request, res *api.Response) error {
	a := new(auth.UserAuth)
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		// API - 注册
		// -> username, password
		// <- token, user
		// !! set-token
		Action(rest.POST | rest.GET).
		Chain(rest.ParamCheck(map[string]*rest.PCC{
			"username": rest.PccMust,
			"password": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, a)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			userToken, err := s.AuthUserClient.(auth.UserService).Register(ctx, a)
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			s.FreshJWT(userToken.Token)
			return s.Success(userToken)
		}).
		Done().
		// Finish
		Final()
}

// Info API
func (*User) Info(ctx context.Context, req *api.Request, res *api.Response) error {
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		// API - 自身信息
		// <- token, user
		// !! set-token
		Action(rest.POST | rest.GET).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			// didn't JWTCheck, beack Info API will also check it
			user, err := s.AuthUserClient.(auth.UserService).Info(s.Ctx, &auth.UserToken{
				Token: rest.GetJWT(s.Req),
			})
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			s.FreshJWT(user.Token)
			return s.Success(user.User)
		}).
		Done().
		// Finish
		Final()
}

// Update API change pass or userinfo
func (*User) Update(ctx context.Context, req *api.Request, res *api.Response) error {
	u := new(auth.UserToken)
	u.Auth = new(auth.UserAuth)
	u.User = new(auth.UserInfo)

	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		// API - 更改自身信息
		// -> pass/userinfo (optional, pass have priority)
		// <- token, user
		// !! set-token
		Action(rest.POST | rest.GET).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(rest.ParamAutoLoad(nil, u.User)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			if u.Auth.Password == "" {
				u.User.UUID = s.Token.User.UUID
				u.Auth = nil
			} else {
				u.Auth.ID = s.Token.User.UUID
				u.User = nil
			}
			user, err := s.AuthUserClient.(auth.UserService).Change(ctx, u)
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			s.FreshJWT(user.Token)
			return s.Success(user.User)
		}).
		Done().
		// Finish
		Final()
}
