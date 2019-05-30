package common

import (
	"context"

	//_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/mssql"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"gopkg.in/mgo.v2"
)

// MongoRepo Wrapper
func MongoRepo(service micro.Service) (wrapper server.HandlerWrapper, deferFn func()) {
	var (
		repo *mgo.Session
		e    error
	)
	if repo, e = mgo.Dial(ENV.Repo.Mongo.URL); e != nil {
		log.Fatal(e)
	}
	repo.SetMode(mgo.Monotonic, true)

	return func(fn server.HandlerFunc) server.HandlerFunc {
			return func(ctx context.Context, req server.Request, rsp interface{}) error {
				ctx = context.WithValue(ctx, RepoKey{CKMongoRepo}, repo)
				return fn(ctx, req, rsp)
			}
		}, func() {
			repo.Close()
		}
}

// GetMongoRepo from ctx
func GetMongoRepo(ctx context.Context) (*mgo.Session, bool) {
	c, ok := ctx.Value(RepoKey{CKMongoRepo}).(*mgo.Session)
	return c.Clone(), ok
}

// GormRepo Wrapper
func GormRepo(service micro.Service, repo *gorm.DB) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, RepoKey{CKGormRepo}, repo)
			return fn(ctx, req, rsp)
		}
	}
}

// GetGormRepo from ctx
func GetGormRepo(ctx context.Context) (*gorm.DB, bool) {
	c, ok := ctx.Value(RepoKey{CKGormRepo}).(*gorm.DB)
	return c, ok
}
