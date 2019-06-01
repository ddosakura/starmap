package handler

import (
	"context"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
)

// Roles Action
func (s *User) Roles(context.Context, *proto.None, *proto.Result) error {
	return nil
}

// Permissions Action
func (s *User) Permissions(context.Context, *proto.None, *proto.Result) error {
	return nil
}
