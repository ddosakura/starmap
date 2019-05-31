package handler

import (
	"context"
	"time"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/jinzhu/gorm"
)

const (
	// TODO: 考虑过期时间
	jwtTerm = time.Hour * 24
)

// User Handler
type User struct{}

// Login Action
func (s *User) Login(ctx context.Context, req *proto.UserAuth, res *proto.UserToken) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	auth := new(proto.UserAuth)
	if err := s.findUser(repo, req.Username, auth); err != nil {
		return raw.ErrUserNotExist
	}
	if auth.Password != req.Password {
		return raw.ErrPassWrong
	}

	// TODO: 获取用户信息
	if token, err := buildUserJWT(nil).sign(jwtTerm); err == nil {
		res.Token = token
		return nil
	}
	return raw.ErrSignJWT
}

func (s *User) findUser(repo *gorm.DB, name string, user *proto.UserAuth) error {
	return repo.First(user, "username = ?", name).Error
}

// Register Action
func (s *User) Register(ctx context.Context, req *proto.UserAuth, res *proto.UserToken) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}
	//repo.Lock()
	//defer repo.Unlock()
	tx := repo.Begin()
	defer tx.Commit()

	// fmt.Println("start")

	auth := new(proto.UserAuth)
	if err := s.findUser(repo, req.Username, auth); err != nil {
		if err == gorm.ErrRecordNotFound {
			u := models.UserAuth{
				Model:    new(common.Model),
				UserAuth: req,
			}
			if err = repo.Create(u).Error; err == nil {
				// TODO: 初始化用户信息
				if token, err := buildUserJWT(nil).sign(jwtTerm); err == nil {
					res.Token = token
					return nil
				}
				return raw.ErrSignJWT
			}
			tx.Rollback()
		}
		return raw.ErrRepoError
	}
	return raw.ErrUserHasExist
}

// Logout Action
func (s *User) Logout(context.Context, *proto.UserToken, *proto.UserToken) error {
	return nil
}

// Info Action
func (s *User) Info(context.Context, *proto.UserToken, *proto.UserToken) error {
	return nil
}

// Change Action
func (s *User) Change(context.Context, *proto.UserToken, *proto.UserToken) error {
	return nil
}

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
