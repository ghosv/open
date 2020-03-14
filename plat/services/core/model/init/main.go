package main

import (
	"github.com/ghosv/open/plat/services/core/initialize"
	model "github.com/ghosv/open/plat/services/core/model"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// TODO: change
const icon = "https://gss0.bdstatic.com/70cFfyinKgQIm2_p8IuM_a/daf/pic/item/80cb39dbb6fd5266c7f1bf4ca418972bd507368b.jpg"

func main() {
	repo := initialize.Repo().Debug()
	defer repo.Close()

	repo.AutoMigrate(
		// &model.User{},
		// &model.App{},

		&model.Org{},
		&model.Group{},
	)

	repo.Where(model.App{Name: "个人中心"}).
		Attrs(model.App{
			// Model: meta.Model{ID: meta.SrvSelf}, // check
			Name:  "个人中心",
			Icon:  icon,
			Intro: "Ghost 个人中心",
			URL:   "https://localhost:3001", // TODO: change
		}).
		FirstOrCreate(&model.App{})
	repo.Where(model.App{Name: "会议室管理"}).
		Attrs(model.App{
			Name:  "会议室管理",
			Icon:  icon,
			Intro: "智能会议室管理系统",
			URL:   "https://localhost:3002", // TODO: change
		}).
		FirstOrCreate(&model.App{})
	//repo.Create(&model.App{
	//	ID:    meta.SrvSelf, // TODO: check
	//	Name:  "个人中心",
	//	Icon:  icon,
	//	Intro: "Ghost 个人中心",
	//	URL:   "https://localhost:3001", // TODO: change
	//})
	//repo.Create(&model.App{
	//	Name:  "会议室管理",
	//	Icon:  icon,
	//	Intro: "智能会议室管理系统",
	//	URL:   "https://localhost:3002", // TODO: change
	//})
}
