package initialize

import (
	"fmt"

	"github.com/ghosv/open/plat/conf"
	model "github.com/ghosv/open/plat/services/core/model"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2/util/log"
)

// Repo init
func Repo() *gorm.DB {
	url := fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.ENV.Repo.MySQL.User,
		conf.ENV.Repo.MySQL.Pass,
		conf.ENV.Repo.MySQL.Host,
		model.RepoName,
	)
	repo, e := gorm.Open("mysql", url)
	if e != nil {
		log.Fatal(e)
	}
	repo.AutoMigrate(
		&model.User{},
		&model.App{},
		&model.AccessScope{},
	)
	return repo
}

// Redis init
func Redis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.ENV.Repo.Redis.URL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
