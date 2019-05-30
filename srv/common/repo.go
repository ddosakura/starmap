package common

import (
	"context"

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
				ctx = context.WithValue(ctx, CKMongoRepo, repo)
				return fn(ctx, req, rsp)
			}
		}, func() {
			repo.Close()
		}
}

// GetMongoRepo from ctx
func GetMongoRepo(ctx context.Context) (*mgo.Session, bool) {
	c, ok := ctx.Value(CKMongoRepo).(*mgo.Session)
	return c.Clone(), ok
}
