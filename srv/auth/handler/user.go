package handler

import (
	"context"

	"gopkg.in/mgo.v2/bson"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/jinzhu/gorm"
)

// Insert Action
func (s *User) Insert(ctx context.Context, req *proto.UserAuth, res *proto.UserToken) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}
	//repo.Lock()
	//defer repo.Unlock()
	tx := repo.Begin()
	defer tx.Commit()
	defer func() {
		e := recover()
		if e != nil {
			tx.Rollback()
		}
	}()

	auth := new(proto.UserAuth)
	var result *gorm.DB
	if result = s.findUser(repo, req.Username, auth); result.Error == nil {
		return raw.ErrUserHasExist
	}
	if !result.RecordNotFound() {
		return raw.ErrRepoError
	}
	u := models.User{
		Model:    new(common.Model),
		UserAuth: req,
	}
	if result = repo.Create(u); result.Error != nil {
		tx.Rollback()
		return raw.ErrRepoError
	}
	mongo, ok := common.GetMongoRepo(ctx)
	if !ok {
		tx.Rollback()
		return raw.ErrRepoNotFound
	}
	user := &proto.UserInfo{
		UUID:     u.ID,
		Nickname: "user-" + u.ID,
		Avatar:   "",
		Motto:    "",
		Phone:    "",
		Email:    "",
		Homepage: "",
	}
	if err := mongo.DB(raw.UserDB).C(raw.UserInfoC).Insert(user); err != nil {
		tx.Rollback()
		return raw.ErrRepoError
	}
	res.User = user
	return nil
}

// Delete Action
func (s *User) Delete(ctx context.Context, req *proto.UserAuth, res *proto.UserToken) error {
	repo, ok := common.GetGormRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	if result := repo.Delete(&models.User{}, "id = ?", req.ID); result.Error != nil {
		if result.RecordNotFound() {
			return raw.ErrUserNotExist
		}
		return raw.ErrRepoError
	} else if result.RowsAffected == 0 {
		return raw.ErrUserNotExist
	}

	return nil
}

// Select Action
func (s *User) Select(ctx context.Context, req *proto.UserAuth, res *proto.UserToken) error {
	if req.ID == "" {
		repo, ok := common.GetGormRepo(ctx)
		if !ok {
			return raw.ErrRepoNotFound
		}
		if r := repo.Where("username = ?", req.Username).First(&models.User{UserAuth: req}); r.Error != nil {
			if r.RecordNotFound() {
				return raw.ErrUserNotExist
			}
			return raw.ErrRepoError
		}
	}
	mongo, ok := common.GetMongoRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}

	user := &proto.UserInfo{}
	c := mongo.DB(raw.UserDB).C(raw.UserInfoC)
	if err := c.Find(&bson.M{"uuid": req.ID}).One(user); err != nil {
		return raw.ErrRepoError
	}

	res.User = user
	return nil
}

// Update Action
func (s *User) Update(ctx context.Context, req *proto.UserToken, res *proto.UserToken) error {
	return nil
}
