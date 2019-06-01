package handler

import (
	"context"

	// proto "github.com/ddosakura/starmap/api/sys/proto"
	api "github.com/micro/go-api/proto"
)

// Role Handler
type Role struct{}

// Entity API
func (*Role) Entity(ctx context.Context, req *api.Request, res *api.Response) error {
	return nil
}
