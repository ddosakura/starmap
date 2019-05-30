package subscriber

import (
	"context"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/go-log/log"
)

// User Sub
type User struct{}

// Handle Auth
func (a *User) Handle(ctx context.Context, msg *proto.UserInfo) error {
	log.Log("Handler Received message: ", msg)
	return nil
}

// Handler for Auth
func Handler(ctx context.Context, msg *proto.UserInfo) error {
	log.Log("Function Received message: ", msg)
	return nil
}
