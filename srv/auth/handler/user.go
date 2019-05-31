package handler

import (
	"context"
	"time"

	"github.com/ddosakura/starmap/srv/auth/models"
	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2/bson"
)

const (
	// TODO: 考虑过期时间、刷新时间
	jwtTerm  = time.Hour * 24
	jwtFresh = time.Hour * 6
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

	user := &proto.UserInfo{}
	repo2, ok := common.GetMongoRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}
	c := repo2.DB(raw.UserDB).C(raw.UserInfoC)
	if err := c.Find(&bson.M{"uuid": auth.ID}).One(user); err != nil {
		return raw.ErrRepoError
	}
	if token, err := common.BuildUserJWT(user).Sign(jwtTerm); err == nil {
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
				repo, ok := common.GetMongoRepo(ctx)
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
				if err := repo.DB(raw.UserDB).C(raw.UserInfoC).Insert(user); err == nil {
					if token, err := common.BuildUserJWT(user).Sign(jwtTerm); err == nil {
						res.Token = token
						return nil
					}
					return raw.ErrSignJWT
				}
			}
			tx.Rollback()
		}
		return raw.ErrRepoError
	}
	return raw.ErrUserHasExist
}

// Info Action
func (s *User) Info(ctx context.Context, req *proto.UserToken, res *proto.UserToken) error {
	repo, ok := common.GetMongoRepo(ctx)
	if !ok {
		return raw.ErrRepoNotFound
	}
	token, err := common.ValidUserJWT(req.Token)
	if err != nil {
		return err
	}
	user := &proto.UserInfo{}
	c := repo.DB(raw.UserDB).C(raw.UserInfoC)
	if err := c.Find(&bson.M{"uuid": token.UserInfo.UUID}).One(user); err != nil {
		return raw.ErrRepoError
	}
	if token, err := common.BuildUserJWT(user).Sign(jwtTerm); err == nil {
		res.Token = token
		return nil
	}
	return raw.ErrSignJWT
}

// Check Action
func (s *User) Check(ctx context.Context, req *proto.UserToken, res *proto.UserToken) error {
	token, err := common.ValidUserJWT(req.Token)
	if err != nil {
		return err
	}
	Token, err := common.FreshJWT(token, jwtFresh, jwtTerm)
	if err != nil {
		return err
	}
	if Token != "" {
		res.Token = Token
	}
	res.User = token.UserInfo
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
