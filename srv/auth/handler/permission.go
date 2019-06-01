package handler

import (
	"context"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
)

// Perm Handler
type Perm struct{}

// Insert API
func (s *Perm) Insert(context.Context, *proto.PermInfo, *proto.PermInfo) error {
	return nil
}

// Delete API
func (s *Perm) Delete(context.Context, *proto.PermInfo, *proto.PermInfo) error {
	return nil
}

// Select API
func (s *Perm) Select(context.Context, *proto.PermInfo, *proto.PermInfo) error {
	return nil
}

// Update API
func (s *Perm) Update(context.Context, *proto.PermInfo, *proto.PermInfo) error {
	return nil
}
