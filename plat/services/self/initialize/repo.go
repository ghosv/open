package initialize

import (
	"fmt"

	"github.com/ghosv/open/plat/conf"
	"github.com/ghosv/open/plat/services/core/model"
	coreModel "github.com/ghosv/open/plat/services/core/model"

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
		coreModel.RepoName,
	)
	repo, e := gorm.Open("mysql", url)
	if e != nil {
		log.Fatal(e)
	}
	repo.AutoMigrate(
		&model.User{},
		&model.App{},
		&model.AccessScope{},
		&model.Org{},
		&model.Group{},
	)
	return repo
}
