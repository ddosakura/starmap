package handler

import (
	"context"
	"fmt"

	"github.com/ddosakura/starmap/api/rest"

	// proto "github.com/ddosakura/starmap/api/sys/proto"
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

	// TODO: User Entity API
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		// API
		Action(rest.POST).
		Chain(rest.PermCheck([]string{"user:insert"}, rest.LogicalAND)).
		Chain(rest.ParamCheck(rest.PCCS{
			"username": rest.PccMust,
			"password": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			userToken, err := s.AuthUserClient.(auth.UserService).Insert(ctx, u.Auth)
			if err != nil {
				return err
			}
			return s.Success(userToken.User)
		}).
		Done().
		// API
		Action(rest.DELETE).
		Chain(rest.PermCheck([]string{"user:delete"}, rest.LogicalAND)).
		Chain(rest.ParamCheck(rest.PCCS{
			"username": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest, u)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.GET).
		Chain(rest.PermCheck([]string{"user:select"}, rest.LogicalAND)).
		Chain(rest.ParamCheck(rest.PCCS{
			"username": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest, u)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.PUT).
		Chain(rest.PermCheck([]string{"user:update"}, rest.LogicalAND)).
		Chain(rest.ParamCheck(rest.PCCS{
			"username": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, u.Auth)).
		Chain(rest.ParamAutoLoad(nil, u.User)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest, u)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// Finish
		Final()
}

// Role Modify API
func (*User) Role(ctx context.Context, req *api.Request, res *api.Response) error {
	// TODO: Role Modify API
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		Chain(rest.PermCheck([]string{"user:role"}, rest.LogicalAND)).
		// API
		Action(rest.POST).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.DELETE).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.GET).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// Finish
		Final()
}
