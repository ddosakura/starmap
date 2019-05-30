package handler

import (
	"context"
	"fmt"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/jinzhu/gorm"
)

// User Handler
type User struct{}

// Login Action
func (s *User) Login(ctx context.Context, req *proto.UserInfo, res *proto.UserInfo) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		res = nil
		return raw.ErrRepoNotFound
	}

	if err := s.findUser(repo, req.Username, res); err != nil {
		res = nil
		return raw.ErrUserNotExist
	}
	if res.Password != req.Password {
		res = nil
		return raw.ErrPassWrong
	}
	res.Password = ""
	return nil
}

func (s *User) findUser(repo *gorm.DB, name string, user *proto.UserInfo) error {
	return repo.First(user, "username = ?", name).Error
}

// Register Action
func (s *User) Register(ctx context.Context, req *proto.UserInfo, res *proto.UserInfo) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		res = nil
		return raw.ErrRepoNotFound
	}
	//repo.Lock()
	//defer repo.Unlock()
	tx := repo.Begin()
	defer tx.Commit()

	fmt.Println("start")

	if err := s.findUser(repo, req.Username, res); err != nil {
		if err == gorm.ErrRecordNotFound {
			err = repo.Create(models.UserInfo{
				new(common.Model),
				req,
				new(models.BeforeCreateHook),
			}).Error
			if err == nil {
				return nil
			}
			tx.Rollback()
		}
		return raw.ErrRepoError
	}
	res = nil
	return raw.ErrUserHasExist
}

// Logout Action
func (s *User) Logout(context.Context, *proto.UserInfo, *proto.UserInfo) error {
	return nil
}

// Info Action
func (s *User) Info(context.Context, *proto.UserInfo, *proto.UserInfo) error {
	return nil
}

// Update Action
func (s *User) Update(context.Context, *proto.UserInfo, *proto.UserInfo) error {
	return nil
}
