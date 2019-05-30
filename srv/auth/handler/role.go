package handler

import (
	"context"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
)

// Role Handler
type Role struct {
}

// Insert API
func (s *Role) Insert(context.Context, *proto.RoleInfo, *proto.RoleInfo) error {
	return nil
}

// Delete API
func (s *Role) Delete(context.Context, *proto.RoleInfo, *proto.RoleInfo) error {
	return nil
}

// Select API
func (s *Role) Select(context.Context, *proto.RoleInfo, *proto.RoleInfo) error {
	return nil
}

// Update API
func (s *Role) Update(context.Context, *proto.RoleInfo, *proto.RoleInfo) error {
	return nil
}
