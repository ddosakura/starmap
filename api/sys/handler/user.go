package handler

import (
	"context"
	"fmt"

	"github.com/ddosakura/starmap/api/rest"
	"github.com/ddosakura/starmap/api/sys/client"

	// proto "github.com/ddosakura/starmap/api/sys/proto"
	api "github.com/micro/go-api/proto"
)

// User Handler
type User struct{}

// Entity API
func (*User) Entity(ctx context.Context, req *api.Request, res *api.Response) error {
	// TODO: User Entity API
	return rest.REST(ctx, req, res).
		Chain(rest.LoadAuthService(client.AuthUserFromContext)).
		Chain(rest.JWTCheck()).
		// API
		Action(rest.POST).
		Chain(rest.PermissionCheck([]string{"user:insert"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.DELETE).
		Chain(rest.PermissionCheck([]string{"user:delete"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.GET).
		Chain(rest.PermissionCheck([]string{"user:select"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.PUT).
		Chain(rest.PermissionCheck([]string{"user:update"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// Finish
		Final()
}
