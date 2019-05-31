package handler

import (
	"context"

	"github.com/ddosakura/starmap/api/auth/client"
	"github.com/ddosakura/starmap/api/auth/raw"
	"github.com/ddosakura/starmap/api/common"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// User Handler
type User struct {
}

// Login API
func (u *User) Login(ctx context.Context, req *api.Request, res *api.Response) error {
	return common.
		REST(ctx, req, res).
		LoadAuthService(client.AuthUserFromContext).
		ACTION(common.POST|common.GET).
		Check("user", false, nil).
		Check("pass", false, nil).
		Do(func(s *common.RESTful) (interface{}, error) {
			userToken, err := s.AuthUserClient.Login(ctx, &auth.UserAuth{
				Username: s.Params["user"].Values[0],
				Password: s.Params["pass"].Values[0],
			})
			if err != nil {
				return nil, common.CleanErrResponse(raw.SrvName, err, errors.BadRequest)
			}
			s.FreshJWT(userToken.Token)
			return userToken, nil
		}).
		Final()
}

// Register API
func (u *User) Register(ctx context.Context, req *api.Request, res *api.Response) error {
	return common.
		REST(ctx, req, res).
		LoadAuthService(client.AuthUserFromContext).
		ACTION(common.POST|common.GET).
		Check("user", false, nil).
		Check("pass", false, nil).
		Do(func(s *common.RESTful) (interface{}, error) {
			userToken, err := s.AuthUserClient.Register(ctx, &auth.UserAuth{
				Username: s.Params["user"].Values[0],
				Password: s.Params["pass"].Values[0],
			})
			if err != nil {
				return nil, common.CleanErrResponse(raw.SrvName, err, errors.BadRequest)
			}
			s.FreshJWT(userToken.Token)
			return userToken, nil
		}).
		Final()
}

// Info API
func (u *User) Info(ctx context.Context, req *api.Request, res *api.Response) error {
	return common.
		REST(ctx, req, res).
		LoadAuthService(client.AuthUserFromContext).
		ACTION(common.POST | common.GET).
		Do(func(s *common.RESTful) (interface{}, error) {
			user, err := s.AuthUserClient.Check(s.Ctx, &auth.UserToken{
				Token: common.GetJWT(s.Req),
			})
			if err != nil {
				return nil, err
			}
			s.FreshJWT(user.Token)
			return user.User, nil
		}).
		Final()
}

// Update API change pass or userinfo
func (u *User) Update(ctx context.Context, req *api.Request, res *api.Response) error {
	return common.
		REST(ctx, req, res).
		LoadAuthService(client.AuthUserFromContext).
		CheckJWT().
		ACTION(common.POST|common.GET).
		Check("pass", false, []string{""}).
		Check("nickname", false, []string{""}).
		Check("avatar", false, []string{""}).
		Check("motto", false, []string{""}).
		Check("phone", false, []string{""}).
		Check("email", false, []string{""}).
		Check("homepage", false, []string{""}).
		Do(func(s *common.RESTful) (interface{}, error) {
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
			user, err := s.AuthUserClient.Change(ctx, user)
			if err != nil {
				return nil, err
			}
			s.FreshJWT(user.Token)
			return user.User, nil
		}).
		Final()
}
