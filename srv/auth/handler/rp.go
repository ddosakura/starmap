package handler

import (
	"context"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
)

// Roles Action
func (s *User) Roles(ctx context.Context, req *proto.None, res *proto.Result) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	var roles []models.Role
	ua := new(proto.UserAuth)
	ua.ID = req.UUID
	user := new(models.User)
	user.UserAuth = ua
	// TODO: gorm's error
	if err := repo.Model(user).Related(&roles, "Roles").Error; err != nil {
		return err
	}

	res.Data = make([]string, 0, len(roles))
	for _, r := range roles {
		// TODO: use roleName after
		res.Data = append(res.Data, r.ID)
	}
	return nil
}

// Permissions Action
func (s *User) Permissions(ctx context.Context, req *proto.None, res *proto.Result) error {
	return nil
}

// Permissions Action
func (s *Role) Permissions(ctx context.Context, req *proto.None, res *proto.Result) error {
	return nil
}
