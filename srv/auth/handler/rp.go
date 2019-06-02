package handler

import (
	"context"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
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
	role := new(models.Role)
	//role.RoleInfo = new(proto.RoleInfo)
	//role.RoleInfo.Name = req.Name
	if result := repo.First(role, "name = ?", req.Name); result.Error != nil {
		if result.RecordNotFound() {
			return raw.ErrRoleNotExist
		}
		return raw.ErrRepoError
	}
	if result := repo.Model(role).Related(&perms, "Perms"); result.Error != nil {
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

	switch req.Modify {
	case proto.M_Add, proto.M_Del:
		user := new(models.User)
		user.UserAuth = new(proto.UserAuth)
		user.UserAuth.ID = req.UUID
		var role models.Role
		if r := repo.First(&role, "name = ?", req.Name); r.Error != nil {
			if r.RecordNotFound() {
				return raw.ErrRoleNotExist
			}
			return raw.ErrRepoError
		}
		if r := repo.First(user, "id = ?", req.UUID); r.Error != nil {
			if r.RecordNotFound() {
				return raw.ErrUserNotExist
			}
			return raw.ErrRepoError
		}
		if req.Modify == proto.M_Add {
			if r := repo.Model(user).Association("Roles").Append(role); r.Error != nil {
				return raw.ErrRepoError
			}
		} else {
			if r := repo.Model(user).Association("Roles").Delete(role); r.Error != nil {
				return raw.ErrRepoError
			}
		}

		// 返回该用户所有角色
		return s.Roles(ctx, &proto.Identity{UUID: req.UUID}, res)
	case proto.M_List:
		var roles []models.Role
		if r := repo.Find(&roles); r.Error != nil {
			return raw.ErrRepoError
		}
		res.Data = make([]string, 0, len(roles))
		for _, r := range roles {
			res.Data = append(res.Data, r.Name)
		}
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

	switch req.Modify {
	case proto.M_Add, proto.M_Del:
		role := new(models.Role)
		role.RoleInfo = new(proto.RoleInfo)
		role.RoleInfo.Name = req.Name
		var perm models.Perm
		if r := repo.First(&perm, "subject = ? AND action = ?", req.Subject, req.Action); r.Error != nil {
			if r.RecordNotFound() {
				return raw.ErrPermNotExist
			}
			return raw.ErrRepoError
		}
		if r := repo.First(role, "name = ?", req.Name); r.Error != nil {
			if r.RecordNotFound() {
				return raw.ErrUserNotExist
			}
			return raw.ErrRepoError
		}
		if req.Modify == proto.M_Add {
			if r := repo.Model(role).Association("Perms").Append(perm); r.Error != nil {
				return raw.ErrRepoError
			}
		} else {
			if r := repo.Model(role).Association("Perms").Delete(perm); r.Error != nil {
				return raw.ErrRepoError
			}
		}

		// 返回该角色所有权限
		return s.Perms(ctx, &proto.Identity{Name: req.Name}, res)
	case proto.M_List:
		var perms []models.Perm
		if r := repo.Find(&perms); r.Error != nil {
			return raw.ErrRepoError
		}
		res.Data = make([]string, 0, len(perms))
		for _, p := range perms {
			res.Data = append(res.Data, p.Subject+":"+p.Action)
		}
	default:
		return raw.ErrUnknowAction
	}
	return nil
}
