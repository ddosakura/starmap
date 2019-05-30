package handler

import (
	"context"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
)

// Permission Handler
type Permission struct{}

// Insert API
func (s *Permission) Insert(context.Context, *proto.PermissionInfo, *proto.PermissionInfo) error {
	return nil
}

// Delete API
func (s *Permission) Delete(context.Context, *proto.PermissionInfo, *proto.PermissionInfo) error {
	return nil
}

// Select API
func (s *Permission) Select(context.Context, *proto.PermissionInfo, *proto.PermissionInfo) error {
	return nil
}

// Update API
func (s *Permission) Update(context.Context, *proto.PermissionInfo, *proto.PermissionInfo) error {
	return nil
}
