package handler

import (
	"context"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
)

// Insert Action
func (s *User) Insert(context.Context, *proto.UserToken, *proto.UserToken) error {
	return nil
}

// Delete Action
func (s *User) Delete(context.Context, *proto.UserToken, *proto.UserToken) error {
	return nil
}

// Select Action
func (s *User) Select(context.Context, *proto.UserToken, *proto.UserToken) error {
	return nil
}

// Update Action
func (s *User) Update(context.Context, *proto.UserToken, *proto.UserToken) error {
	return nil
}
