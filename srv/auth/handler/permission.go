package handler

import (
	"context"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/jinzhu/gorm"
)

// Perm Handler
type Perm struct{}

// Insert API
func (s *Perm) Insert(ctx context.Context, req *proto.PermInfo, res *proto.PermInfo) error {
	if req.Subject == "" || req.Action == "" {
		return raw.ErrPermHasExist
	}
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	if r := repo.Where(pQuery, req.Subject, req.Action).Attrs(req).FirstOrCreate(&models.Perm{PermInfo: res}); r.Error != nil {
		return raw.ErrRepoError
	}

	if req.Detail != res.Detail {
		return raw.ErrPermHasExist
	}
	return nil
}

const pQuery = "subject = ? AND action = ?"

func pSelect(repo *gorm.DB, sub, act string, entity *models.Perm) error {
	if r := repo.Where(pQuery, sub, act).First(entity); r.Error != nil {
		if r.RecordNotFound() {
			return raw.ErrPermNotExist
		}
		return raw.ErrRepoError
	}
	return nil
}

// Delete API
func (s *Perm) Delete(ctx context.Context, req *proto.PermInfo, res *proto.PermInfo) error {
	if req.Subject == "" || req.Action == "" {
		return raw.ErrPermHasExist
	}
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	p := &models.Perm{PermInfo: res}
	if e := pSelect(repo, req.Subject, req.Action, p); e != nil {
		return e
	}
	return repo.Delete(p).Error
}

// Select API
func (s *Perm) Select(ctx context.Context, req *proto.PermInfo, res *proto.PermInfo) error {
	if req.Subject == "" || req.Action == "" {
		return raw.ErrPermHasExist
	}
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	p := &models.Perm{PermInfo: res}
	return pSelect(repo, req.Subject, req.Action, p)
}

// Update API
func (s *Perm) Update(ctx context.Context, req *proto.PermWrapper, res *proto.PermInfo) error {
	if req.Subject == "" || req.Action == "" {
		return raw.ErrPermHasExist
	}
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	p := &models.Perm{PermInfo: res}
	if e := pSelect(repo, req.Subject, req.Action, p); e != nil {
		return e
	}
	if req.Perm.Subject != "" {
		res.Subject = req.Perm.Subject
	}
	if req.Perm.Action != "" {
		res.Action = req.Perm.Action
	}
	if req.Perm.Detail != "" {
		res.Detail = req.Perm.Detail
	}
	if r := repo.Save(p); r.Error != nil {
		return raw.ErrRepoError
	}
	return nil
}
