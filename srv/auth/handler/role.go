package handler

import (
	"context"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/jinzhu/gorm"
	"github.com/kr/pretty"
)

// Role Handler
type Role struct {
}

// Insert API
func (s *Role) Insert(ctx context.Context, req *proto.RoleInfo, res *proto.RoleInfo) error {
	if req.Name == "" {
		return raw.ErrRoleHasExist
	}
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	pretty.Println(req)
	if r := repo.Debug().Where("name = ?", req.Name).Attrs(req).FirstOrCreate(&models.Role{RoleInfo: res}); r.Error != nil {
		return raw.ErrRepoError
	}

	if req.Detail != res.Detail {
		return raw.ErrRoleHasExist
	}
	return nil
}

// Delete API
func (s *Role) Delete(ctx context.Context, req *proto.RoleInfo, res *proto.RoleInfo) error {
	if req.Name == "" {
		return raw.ErrRoleHasExist
	}
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}
	r := &models.Role{RoleInfo: res}
	if e := selectQ(repo, req.Name, r); e != nil {
		return e
	}
	return repo.Debug().Delete(r).Error
}

func selectQ(repo *gorm.DB, name string, entity *models.Role) error {
	if r := repo.Debug().Where("name = ?", name).First(entity); r.Error != nil {
		if r.RecordNotFound() {
			return raw.ErrRoleNotExist
		}
		return raw.ErrRepoError
	}
	return nil
}

// Select API
func (s *Role) Select(ctx context.Context, req *proto.RoleInfo, res *proto.RoleInfo) error {
	if req.Name == "" {
		return raw.ErrRoleHasExist
	}
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	return selectQ(repo, req.Name, &models.Role{RoleInfo: res})
}

// Update API
func (s *Role) Update(ctx context.Context, req *proto.RoleWrapper, res *proto.RoleInfo) error {
	if req.Name == "" {
		return raw.ErrRoleHasExist
	}
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	role := &models.Role{RoleInfo: res}
	if e := selectQ(repo, req.Name, role); e != nil {
		return e
	}
	if req.Role.Name != "" {
		res.Name = req.Role.Name
	}
	if req.Role.Detail != "" {
		res.Detail = req.Role.Detail
	}
	if r := repo.Debug().Save(role); r.Error != nil {
		return raw.ErrRepoError
	}
	return nil
}
