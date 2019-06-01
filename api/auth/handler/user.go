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
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		// API
		Action(rest.POST | rest.GET).
		Chain(rest.ParamCheck(map[string]*rest.PCC{
			"user": &rest.PCC{Must: true},
			"pass": &rest.PCC{Must: true},
		})).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			userToken, err := s.AuthUserClient.(auth.UserService).Login(ctx, &auth.UserAuth{
				Username: s.Params["user"].Values[0],
				Password: s.Params["pass"].Values[0],
			})
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
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		// API
		Action(rest.POST | rest.GET).
		Chain(rest.ParamCheck(rest.PCCS{
			"user": &rest.PCC{Must: true},
			"pass": &rest.PCC{Must: true},
		})).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			userToken, err := s.AuthUserClient.(auth.UserService).Register(ctx, &auth.UserAuth{
				Username: s.Params["user"].Values[0],
				Password: s.Params["pass"].Values[0],
			})
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
		// API
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
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		// API
		Action(rest.POST | rest.GET).
		Chain(rest.ParamCheck(rest.PCCS{
			"pass":     rest.PccEmptyStr,
			"nickname": rest.PccEmptyStr,
			"avatar":   rest.PccEmptyStr,
			"motto":    rest.PccEmptyStr,
			"phone":    rest.PccEmptyStr,
			"email":    rest.PccEmptyStr,
			"homepage": rest.PccEmptyStr,
		})).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			user := &auth.UserToken{}
			if s.Params["pass"].Values[0] == "" {
				user.User = &auth.UserInfo{
					UUID:     s.Token.User.UUID,
					Nickname: s.Params["nickname"].Values[0],
					Avatar:   s.Params["avatar"].Values[0],
					Motto:    s.Params["motto"].Values[0],
					Phone:    s.Params["phone"].Values[0],
					Email:    s.Params["email"].Values[0],
					Homepage: s.Params["homepage"].Values[0],
				}
			} else {
				user.Auth = &auth.UserAuth{
					ID:       s.Token.User.UUID,
					Password: s.Params["pass"].Values[0],
				}
			}
			user, err := s.AuthUserClient.(auth.UserService).Change(ctx, user)
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
