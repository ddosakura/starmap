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
	authUserClient, ok := client.AuthUserFromContext(ctx)
	if !ok {
		return errors.InternalServerError(raw.SrvName, "auth client not found")
	}

	return common.
		REST(ctx, req, res).
		ACTION(common.POST|common.GET).
		Check("user", false, nil).
		Check("pass", false, nil).
		Do(func(s *common.RESTful) (interface{}, error) {
			userInfo, err := authUserClient.Login(ctx, &auth.UserInfo{
				Username: s.Params["user"].Values[0],
				Password: s.Params["pass"].Values[0],
			})
			if err != nil {
				return nil, errors.BadRequest(raw.SrvName, err.Error())
			}
			return userInfo, nil
		}).
		Final()
}

// Register API
func (u *User) Register(ctx context.Context, req *api.Request, res *api.Response) error {
	authUserClient, ok := client.AuthUserFromContext(ctx)
	if !ok {
		return errors.InternalServerError(raw.SrvName, "auth client not found")
	}

	return common.
		REST(ctx, req, res).
		ACTION(common.POST|common.GET).
		Check("user", false, nil).
		Check("pass", false, nil).
		Do(func(s *common.RESTful) (interface{}, error) {
			_, err := authUserClient.Register(ctx, &auth.UserInfo{
				Username: s.Params["user"].Values[0],
				Password: s.Params["pass"].Values[0],
			})
			if err != nil {
				return nil, errors.BadRequest(raw.SrvName, err.Error())
			}
			return nil, nil
		}).
		Final()
}

// Logout API
func (u *User) Logout(context.Context, *api.Request, *api.Response) error {
	return nil
}

// Info API
func (u *User) Info(context.Context, *api.Request, *api.Response) error {
	return nil
}

// Update API
func (u *User) Update(context.Context, *api.Request, *api.Response) error {
	return nil
}
