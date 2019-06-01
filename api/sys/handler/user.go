package handler

import (
	"context"
	"fmt"

	"github.com/ddosakura/starmap/api/common"
	"github.com/ddosakura/starmap/api/sys/client"

	// proto "github.com/ddosakura/starmap/api/sys/proto"
	api "github.com/micro/go-api/proto"
)

// User Handler
type User struct{}

// Entity API
func (*User) Entity(ctx context.Context, req *api.Request, res *api.Response) error {
	return common.
		REST(ctx, req, res).
		LoadAuthService(client.AuthUserFromContext).
		CheckJWT().
		// Role([]string{"admin"}, common.LogicalAND).
		ACTION(common.POST).
		Permission([]string{"user:insert"}, common.LogicalAND).
		Do(func(s *common.RESTful) (interface{}, error) {
			fmt.Println("M", s.Rest)
			return nil, nil
		}).
		ACTION(common.DELETE).
		Permission([]string{"user:delete"}, common.LogicalAND).
		Do(func(s *common.RESTful) (interface{}, error) {
			fmt.Println("M", s.Rest)
			return nil, nil
		}).
		ACTION(common.GET).
		Permission([]string{"user:select"}, common.LogicalAND).
		Do(func(s *common.RESTful) (interface{}, error) {
			fmt.Println("M", s.Rest)
			return nil, nil
		}).
		ACTION(common.PUT).
		Permission([]string{"user:update"}, common.LogicalAND).
		Do(func(s *common.RESTful) (interface{}, error) {
			fmt.Println("M", s.Rest)
			return nil, nil
		}).
		Final()
}
