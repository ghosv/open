package conf

import (
	"context"

	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/log"
	"gopkg.in/mgo.v2"

	"github.com/timest/env"
)

type config struct {
	Repo struct {
		Mongo struct {
			URL string `env:"URL" default:"root:123456@127.0.0.1:27017"`
		} `env:"MONGO"`
		Redis struct {
			URL string `env:"URL" default:"127.0.0.1:6379"`
		} `env:"REDIS"`
		MySQL struct {
			Host string `env:"HOST" default:"127.0.0.1:3306"`
			User string `env:"USER" default:"root"`
			Pass string `env:"PASS" default:"123456"`
		} `env:"MYSQL"`
	} `env:"REPO"`
	Key struct {
		JWT string `env:"JWT" default:"jwt_key"`
	} `env:"KEY"`
}

func (c *config) KeyJWT() []byte {
	return []byte(c.Key.JWT)
}

// config(s)
var (
	ENV = new(config)
)

func init() {
	if err := env.Fill(ENV); err != nil {
		log.Fatal(err)
	}
}

// === tmp ===

type keyRepo struct {
	t string
}

// MongoRepo Wrapper
func MongoRepo() (wrapper server.HandlerWrapper, deferFn func()) {
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
				ctx = context.WithValue(ctx, keyRepo{"mgo"}, repo)
				return fn(ctx, req, rsp)
			}
		}, func() {
			repo.Close()
		}
}

// GetMongoRepo from ctx
func GetMongoRepo(ctx context.Context) (*mgo.Session, bool) {
	c, ok := ctx.Value(keyRepo{"mgo"}).(*mgo.Session)
	return c.Clone(), ok
}
