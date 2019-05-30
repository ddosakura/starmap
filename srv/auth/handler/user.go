package handler

import (
	"context"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User Handler
type User struct{}

// Login Action
func (s *User) Login(ctx context.Context, req *proto.UserInfo, res *proto.UserInfo) error {
	repo, ok := common.GetMongoRepo(ctx)
	if !ok {
		res = nil
		return raw.ErrRepoNotFound
	}
	db := repo.DB(raw.AuthDB)

	if err := s.findUser(db.C("userinfo"), req.Username, res); err != nil {
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

func (s *User) findUser(repo *mgo.Collection, name string, user *proto.UserInfo) error {
	return repo.Find(bson.M{"username": name}).One(user)
}

// Register Action
func (s *User) Register(ctx context.Context, req *proto.UserInfo, res *proto.UserInfo) error {
	res = nil
	repo, ok := common.GetMongoRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}
	db := repo.DB(raw.AuthDB)

	c := db.C("userinfo")
	if err := s.findUser(c, req.Username, res); err != nil {
		if err == mgo.ErrNotFound {
			err = c.Insert(req)
			if err == nil {
				return nil
			}
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
