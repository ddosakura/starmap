package handler

import (
	"context"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/kr/pretty"
)

// Roles Action
func (s *User) Roles(ctx context.Context, req *proto.Identity, res *proto.Result) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	var roles []models.Role
	user := new(models.User)
	user.UserAuth = new(proto.UserAuth)
	user.UserAuth.ID = req.UUID
	if result := repo.Model(user).Related(&roles, "Roles"); result.Error != nil {
		return raw.ErrRepoError
	}

	res.Data = make([]string, 0, len(roles))
	for _, r := range roles {
		res.Data = append(res.Data, r.Name)
	}
	return nil
}

// Perms Action
func (s *User) Perms(ctx context.Context, req *proto.Identity, res *proto.Result) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	var roles []models.Role
	var perms []models.Perm
	user := new(models.User)
	user.UserAuth = new(proto.UserAuth)
	user.UserAuth.ID = req.UUID
	if result := repo.Model(user).Related(&roles, "Roles"); result.Error != nil {
		return raw.ErrRepoError
	}
	if result := repo.Model(&roles).Related(&perms, "Perms"); result.Error != nil {
		return raw.ErrRepoError
	}

	res.Data = make([]string, 0, len(perms))
	for _, p := range perms {
		res.Data = append(res.Data, p.Subject+":"+p.Action)
	}
	return nil
}

// Perms Action
func (s *Role) Perms(ctx context.Context, req *proto.Identity, res *proto.Result) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	var perms []models.Perm
	rule := new(models.Role)
	//rule.RoleInfo = new(proto.RoleInfo)
	//rule.RoleInfo.Name = req.Name
	if result := repo.First(rule, "name = ?", req.Name); result.Error != nil {
		if result.RecordNotFound() {
			return raw.ErrRoleNotExist
		}
		return raw.ErrRepoError
	}
	if result := repo.Model(rule).Related(&perms, "Perms"); result.Error != nil {
		return raw.ErrRepoError
	}

	res.Data = make([]string, 0, len(perms))
	for _, p := range perms {
		res.Data = append(res.Data, p.Subject+":"+p.Action)
	}
	return nil
}

// Role modify
func (s *User) Role(ctx context.Context, req *proto.Modification, res *proto.Result) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}
	pretty.Println(req)

	var roles []models.Role
	if r := repo.Find(&roles); r.Error != nil {
		return raw.ErrRepoError
	}
	l := make([]string, 0, len(roles))
	for _, r := range roles {
		l = append(l, r.Name)
	}

	switch req.Modify {
	case proto.M_Add:
		result := new(proto.Result)
		if e := s.Roles(ctx, &proto.Identity{UUID: req.UUID}, result); e != nil {
			return e
		}
		if !common.ItemInList(req.Name, l) {
			return raw.ErrRoleNotExist
		}
		// TODO:
	case proto.M_Del:
	case proto.M_List:
		res.Data = l
	default:
		return raw.ErrUnknowAction
	}
	return nil
}

// Perm modify
func (s *Role) Perm(ctx context.Context, req *proto.Modification, res *proto.Result) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}
	pretty.Println(req)

	var perms []models.Perm
	if r := repo.Find(&perms); r.Error != nil {
		return raw.ErrRepoError
	}
	l := make([]string, 0, len(perms))
	for _, p := range perms {
		l = append(l, p.Subject+":"+p.Action)
	}

	switch req.Modify {
	case proto.M_Add:
	case proto.M_Del:
	case proto.M_List:
		res.Data = l
	default:
		return raw.ErrUnknowAction
	}
	return nil
}
