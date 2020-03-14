package handler

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/ghosv/open/meta"
	model "github.com/ghosv/open/plat/services/core/model"
	pb "github.com/ghosv/open/plat/services/core/proto"
	"github.com/ghosv/open/plat/utils"
)

// App Handler
type App struct{}

// List of App
func (s *App) List(ctx context.Context, req *pb.Identity, res *pb.AppList) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)

	var apps []model.App
	user := new(model.User)
	user.ID = req.UUID // TODO: 从 ctx 获取?
	// TODO: 获取时忽略 self
	if result := repo.Model(user).Related(&apps, "Apps"); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	app := &model.App{
		Model: meta.Model{ID: meta.SrvSelf},
	}
	if result := repo.First(app); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	// TODO: owner
	// scopes
	// accessList

	res.Data = make([]*pb.AppInfo, 0, len(apps)+1)
	res.Data = append(res.Data, &pb.AppInfo{
		ID:      app.ID,
		Name:    app.Name,
		Icon:    app.Icon,
		Intro:   app.Intro,
		URL:     app.URL,
		OwnerID: app.OwnerID,
	})
	for _, v := range apps {
		res.Data = append(res.Data, &pb.AppInfo{
			ID:      v.ID,
			Name:    v.Name,
			Icon:    v.Icon,
			Intro:   v.Intro,
			URL:     v.URL,
			OwnerID: v.OwnerID,
		})
	}
	return nil
}

// FindByID App
func (s *App) FindByID(ctx context.Context, req *pb.AppInfo, res *pb.AppInfo) error {
	app := &model.App{
		Model: meta.Model{ID: req.ID},
	}

	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.First(app); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	res.ID = app.ID
	res.Name = app.Name
	res.Icon = app.Icon
	res.Intro = app.Intro
	res.URL = app.URL
	// TODO: 确认是否需要
	// res.Scopes = strings.Split(app.Scopes, "|")
	return nil
}
