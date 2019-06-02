package handler

import (
	"context"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
)

// Perm Handler
type Perm struct{}

// Insert API
func (s *Perm) Insert(ctx context.Context, req *proto.PermInfo, res *proto.PermInfo) error {
	return nil
}

// Delete API
func (s *Perm) Delete(ctx context.Context, req *proto.PermInfo, res *proto.PermInfo) error {
	return nil
}

// Select API
func (s *Perm) Select(ctx context.Context, req *proto.PermInfo, res *proto.PermInfo) error {
	return nil
}

// Update API
func (s *Perm) Update(ctx context.Context, req *proto.PermWrapper, res *proto.PermInfo) error {
	return nil
}
