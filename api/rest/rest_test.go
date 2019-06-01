package rest

import (
	"context"
	"fmt"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
)

func Example() {
	handle := func(ctx context.Context, s *Flow) error {
		fmt.Println("M", s.Rest)
		return s.Success(fmt.Sprintf("M %v", s.Rest))
	}

	Handle := func(ctx context.Context, req *api.Request, res *api.Response) error {
		return REST(ctx, req, res).
			Chain(LoadAuthService(func(ctx context.Context) (auth.UserService, bool) { return nil, false })).
			Chain(JWTCheck()).
			Chain(RoleCheck([]string{"admin"}, LogicalAND)).
			// API
			Action(POST).
			Chain(PermCheck([]string{"user:insert"}, LogicalAND)).
			Chain(handle).
			Done().
			// API
			Action(DELETE).
			Chain(PermCheck([]string{"user:delete"}, LogicalAND)).
			Chain(handle).
			Done().
			// API
			Action(GET).
			Chain(PermCheck([]string{"user:select"}, LogicalAND)).
			Chain(handle).
			Done().
			// API
			Action(PUT).
			Chain(PermCheck([]string{"user:update"}, LogicalAND)).
			Chain(handle).
			Done().
			// Finsh
			Final()
	}

	Handle(context.Background(),
		&api.Request{},
		&api.Response{},
	)
}
