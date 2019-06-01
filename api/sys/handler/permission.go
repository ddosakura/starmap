package handler

import (
	"context"
	"fmt"

	// proto "github.com/ddosakura/starmap/api/sys/proto"
	"github.com/ddosakura/starmap/api/rest"
	"github.com/ddosakura/starmap/api/sys/client"
	api "github.com/micro/go-api/proto"
)

// Permission Handler
type Permission struct{}

// Entity API
func (*Permission) Entity(ctx context.Context, req *api.Request, res *api.Response) error {
	// TODO: Permission Entity API
	return rest.REST(ctx, req, res).
		Chain(rest.LoadAuthService(client.AuthUserFromContext)).
		Chain(rest.JWTCheck()).
		// API
		Action(rest.POST).
		Chain(rest.PermissionCheck([]string{"permission:insert"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.DELETE).
		Chain(rest.PermissionCheck([]string{"permission:delete"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.GET).
		Chain(rest.PermissionCheck([]string{"permission:select"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.PUT).
		Chain(rest.PermissionCheck([]string{"permission:update"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// Finish
		Final()
}
