package handler

import (
	"context"

	// proto "github.com/ddosakura/starmap/api/sys/proto"
	api "github.com/micro/go-api/proto"
)

// Permission Handler
type Permission struct{}

// Entity API
func (*Permission) Entity(ctx context.Context, req *api.Request, res *api.Response) error {
	return nil
}
