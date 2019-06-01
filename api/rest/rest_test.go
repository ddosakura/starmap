package rest

import (
	"context"
	"fmt"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
)

func Example() {
	handle := func(s *Flow) error {
		fmt.Println("M", s.Rest)
		return s.Success(fmt.Sprintf("M %v", s.Rest))
	}

	Handle := func(ctx context.Context, req *api.Request, res *api.Response) error {
		return REST(ctx, req, res).
			// LoadAuthService(client.AuthUserFromContext).
			Chain(LoadAuthService(func(ctx context.Context) (auth.UserService, bool) { return nil, false })).
			// CheckJWT().
			Chain(JWTCheck()).
			// Role([]string{"admin"}, common.LogicalAND).
			Chain(RoleCheck([]string{"admin"}, LogicalAND)).
			Action(POST).
			// Permission([]string{"user:insert"}, common.LogicalAND).
			Chain(handle).
			Done().
			Action(DELETE).
			// Permission([]string{"user:delete"}, common.LogicalAND).
			Chain(handle).
			Done().
			Action(GET).
			// Permission([]string{"user:select"}, common.LogicalAND).
			Chain(handle).
			Done().
			Action(PUT).
			// Permission([]string{"user:update"}, common.LogicalAND).
			Chain(handle).
			Done().
			Final()
	}

	Handle(context.Background(),
		&api.Request{},
		&api.Response{},
	)
}
