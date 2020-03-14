package handler

import (
	"context"

	"github.com/ghosv/open/meta"
	metaPb "github.com/ghosv/open/meta/proto"
	"github.com/ghosv/open/plat/services/core/model"
	pb "github.com/ghosv/open/plat/services/self/proto"
	"github.com/ghosv/open/plat/utils"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// App Handler
type App struct{}

// === Query ===

// BatchFind App
func (h *App) BatchFind(ctx context.Context, req *pb.BatchID, res *pb.AppMap) error {
	var apps []model.App
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.Where(req.UUID).Find(&apps); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	data := make(map[string]*pb.AppInfo, len(apps))
	for _, v := range apps {
		d, e := h.loadData(repo, &v)
		if e != nil {
			return e
		}
		data[v.ID] = d
	}
	res.Data = data
	return nil
}

func mactoPbac(o []model.AccessScope) []*pb.AccessScope {
	list := make([]*pb.AccessScope, 0, len(o))
	for _, v := range o {
		list = append(list, &pb.AccessScope{
			AppID:  v.AppID,
			Name:   v.Name,
			Detail: v.Detail,
		})
	}
	return list
}

func (h *App) loadData(repo *gorm.DB, v *model.App) (*pb.AppInfo, error) {
	var scopes, accessList []model.AccessScope
	if result := repo.
		Model(&model.App{Model: meta.Model{ID: v.ID}}).
		Related(&scopes, "Scopes"); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}
	if result := repo.
		Model(&model.App{Model: meta.Model{ID: v.ID}}).
		Related(&accessList, "AccessList"); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}

	var managers, developers, users []string
	if result := repo.
		Table("app_managers").
		Where("app_id = ?", v.ID).
		Pluck("user_id", &managers); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}
	if result := repo.
		Table("app_developers").
		Where("app_id = ?", v.ID).
		Pluck("user_id", &developers); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}
	if result := repo.
		Table("app_users").
		Where("app_id = ?", v.ID).
		Pluck("user_id", &users); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}

	return &pb.AppInfo{
		ID:         v.ID,
		Name:       v.Name,
		Icon:       v.Icon,
		Intro:      v.Intro,
		URL:        v.URL,
		OwnerID:    v.OwnerID,
		Scopes:     mactoPbac(scopes),
		AccessList: mactoPbac(accessList),

		Key:        v.Key,
		Secret:     v.Secret,
		Managers:   managers,
		Developers: developers,
		Users:      users,
	}, nil
}

// Search App
func (h *App) Search(ctx context.Context, req *pb.SearchForm, res *pb.AppList) error {
	word, page, size := req.Word, req.Page, req.Size
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	s := repo.Model(&model.App{}).Where("name LIKE ?", "%"+word+"%")
	if result := s.Count(&res.Total); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	var apps []model.App
	if result := s.Limit(size).Offset((page - 1) * size).
		Find(&apps); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	res.List = make([]*pb.AppInfo, 0, len(apps))
	for _, v := range apps {
		d, e := h.loadData(repo, &v)
		if e != nil {
			return e
		}
		res.List = append(res.List, d)
	}
	return nil
}

// === Mutation ===

func (h *App) check(repo *gorm.DB, req *pb.AppModify) (*model.App, error) {
	app := model.App{
		Model: meta.Model{
			ID: req.ID,
		},
	}
	if result := repo.First(&app); result.Error != nil {
		// return nil, utils.RepoErrorFilter(result)
		return nil, meta.ErrAccessDenied
	}
	if app.OwnerID != req.Visitor.UUID {
		return nil, meta.ErrAccessDenied
	}
	return &app, nil
}

func (h *App) seeUpdate(repo *gorm.DB, req *pb.AppModify, res *pb.AppInfo) error {
	app := model.App{
		Model: meta.Model{
			ID: req.ID,
		},
	}
	if result := repo.Where(&app).First(&app); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	res0, err := h.loadData(repo, &app)
	if err != nil {
		return err
	}
	*res = *res0
	return nil
}

// Create App
func (h *App) Create(ctx context.Context, req *pb.AppModify, res *pb.AppInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	app := &model.App{
		Name:  req.Name,
		Icon:  req.Icon,
		Intro: req.Intro,
		URL:   req.URL,

		OwnerID: req.Visitor.UUID,
	}
	if result := repo.Create(app); result.Error != nil {
		return meta.ErrAppHasExist
	}
	res0, err := h.loadData(repo, app)
	if err != nil {
		return err
	}
	*res = *res0
	return nil
}

// Delete App
func (h *App) Delete(ctx context.Context, req *pb.AppModify, res *metaPb.None) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if _, err := h.check(repo, req); err != nil {
		return err
	}

	if result := repo.Delete(&model.App{
		Model: meta.Model{
			ID: req.ID,
		},
	}); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	return nil
}

// Update App
func (h *App) Update(ctx context.Context, req *pb.AppModify, res *pb.AppInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	app, err := h.check(repo, req)
	if err != nil {
		return err
	}

	if req.Scope == nil {
		app.Name = req.Name
		app.Icon = req.Icon
		app.Intro = req.Intro
		app.URL = req.URL
		if req.ResetSecret {
			app.Secret = uuid.NewV4().String()
		}

		if result := repo.Model(app).Updates(app); result.Error != nil {
			return utils.RepoErrorFilter(result)
		}

		return h.seeUpdate(repo, req, res)
	}

	switch req.ModifyScopeType {
	case pb.ModifyScopeType_AddScope:
		scope := req.ID + ":" + req.Scope.Name
		if result := repo.Create(&model.AccessScope{
			AppID:  req.ID,
			Name:   scope,
			Detail: req.Scope.Detail,
		}); result.Error != nil {
			return utils.RepoErrorFilter(result)
		}
	case pb.ModifyScopeType_DelScope:
		scope := req.ID + ":" + req.Scope.Name
		if result := repo.
			Where("name = ?", scope).
			Unscoped().
			Delete(&model.AccessScope{}); result.Error != nil {
			return utils.RepoErrorFilter(result)
		}
	case pb.ModifyScopeType_AddAccess:
		scope := model.AccessScope{
			Name: req.Scope.Name,
		}
		if result := repo.
			Where(&scope).
			First(&scope); result.Error != nil {
			return utils.RepoErrorFilter(result)
		}
		if result := repo.
			Model(app).
			Association("AccessList").
			Append(scope); result.Error != nil {
			return result.Error
		}
	case pb.ModifyScopeType_DelAccess:
		scope := model.AccessScope{
			Name: req.Scope.Name,
		}
		if result := repo.
			Where(&scope).
			First(&scope); result.Error != nil {
			return utils.RepoErrorFilter(result)
		}
		if result := repo.
			Model(app).
			Association("AccessList").
			Delete(scope); result.Error != nil {
			return result.Error
		}
	default:
		return meta.ErrUnsupportedActionType
	}

	return h.seeUpdate(repo, req, res)
}

// === Other ===

// Mark App
func (h *App) Mark(ctx context.Context, req *pb.AppModify, res *pb.AppInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.
		Model(&model.App{
			Model: meta.Model{
				ID: req.ID,
			},
		}).
		Association("Users").
		Append(model.User{
			Model: meta.Model{
				ID: req.Visitor.UUID,
			},
		}); result.Error != nil {
		return result.Error
	}

	return h.seeUpdate(repo, req, res)
}

// Unmark App
func (h *App) Unmark(ctx context.Context, req *pb.AppModify, res *pb.AppInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.
		Model(&model.App{
			Model: meta.Model{
				ID: req.ID,
			},
		}).
		Association("Users").
		Delete(model.User{
			Model: meta.Model{
				ID: req.Visitor.UUID,
			},
		}); result.Error != nil {
		return result.Error
	}

	return h.seeUpdate(repo, req, res)
}
