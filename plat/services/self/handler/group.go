package handler

import (
	"context"

	"github.com/ghosv/open/meta"
	metaPb "github.com/ghosv/open/meta/proto"
	"github.com/ghosv/open/plat/services/core/model"
	pb "github.com/ghosv/open/plat/services/self/proto"
	"github.com/ghosv/open/plat/utils"
	"github.com/jinzhu/gorm"
)

// Group Handler
type Group struct{}

// === Query ===

// BatchFind Group
func (h *Group) BatchFind(ctx context.Context, req *pb.BatchID, res *pb.GroupMap) error {
	var groups []model.Group
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.Where(req.UUID).Find(&groups); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	data := make(map[string]*pb.GroupInfo, len(groups))
	for _, v := range groups {
		d, e := h.loadData(repo, &v)
		if e != nil {
			return e
		}
		data[v.ID] = d
	}
	res.Data = data
	return nil
}

func (h *Group) loadData(repo *gorm.DB, v *model.Group) (*pb.GroupInfo, error) {
	var users []string
	if result := repo.
		Table("group_users").
		Where("group_id = ?", v.ID).
		Pluck("user_id", &users); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}
	return &pb.GroupInfo{
		ID:       v.ID,
		Name:     v.Name,
		Icon:     v.Icon,
		Detail:   v.Detail,
		MasterID: v.MasterID,
		Users:    users,
	}, nil
}

// Search Group
func (h *Group) Search(ctx context.Context, req *pb.SearchForm, res *pb.GroupList) error {
	word, page, size := req.Word, req.Page, req.Size
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	s := repo.
		Model(&model.Group{}).
		Where("name LIKE ?", "%"+word+"%")
	if result := s.Count(&res.Total); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	var groups []model.Group
	if result := s.Limit(size).Offset((page - 1) * size).
		Find(&groups); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	res.List = make([]*pb.GroupInfo, 0, len(groups))
	for _, v := range groups {
		d, e := h.loadData(repo, &v)
		if e != nil {
			return e
		}
		res.List = append(res.List, d)
	}
	return nil
}

// === Mutation ===

func (h *Group) check(repo *gorm.DB, visitor, id string) (*model.Group, error) {
	if id == "" {
		// return nil, meta.ErrGroupNotExist
		return nil, meta.ErrAccessDenied
	}
	data := model.Group{
		Model: meta.Model{
			ID: id,
		},
	}
	if result := repo.First(&data); result.Error != nil {
		// return nil, utils.RepoErrorFilter(result)
		return nil, meta.ErrAccessDenied
	}
	if data.MasterID == visitor {
		return &data, nil
	}
	return nil, meta.ErrAccessDenied
}

func (h *Group) seeUpdate(repo *gorm.DB, req *pb.GroupModify, res *pb.GroupInfo) error {
	data := model.Group{
		Model: meta.Model{
			ID: req.Info.ID,
		},
	}
	if result := repo.Where(&data).First(&data); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	res0, err := h.loadData(repo, &data)
	if err != nil {
		return err
	}
	*res = *res0
	return nil
}

// Create Group
func (h *Group) Create(ctx context.Context, req *pb.GroupModify, res *pb.GroupInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)

	data := &model.Group{
		Name:     req.Info.Name,
		Icon:     req.Info.Icon,
		Detail:   req.Info.Detail,
		MasterID: req.Visitor.UUID,
	}
	if result := repo.Create(data); result.Error != nil {
		return meta.ErrGroupHasExist
	}
	res0, err := h.loadData(repo, data)
	if err != nil {
		return err
	}
	*res = *res0
	return nil
}

// Delete Group
func (h *Group) Delete(ctx context.Context, req *pb.GroupModify, res *metaPb.None) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if _, err := h.check(repo, req.Visitor.UUID, req.Info.ID); err != nil {
		return err
	}

	if result := repo.Delete(&model.Group{
		Model: meta.Model{
			ID: req.Info.ID,
		},
	}); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	return nil
}

// Update Group
func (h *Group) Update(ctx context.Context, req *pb.GroupModify, res *pb.GroupInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	data, err := h.check(repo, req.Visitor.UUID, req.Info.ID)
	if err != nil {
		return err
	}

	if result := repo.First(data); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	data.Name = req.Info.Name
	data.Icon = req.Info.Icon
	data.Detail = req.Info.Detail
	data.MasterID = req.Info.MasterID

	if result := repo.Model(data).Updates(data); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	if len(req.Info.Users) > 0 {
		if result := repo.
			Model(data).
			Association("Users").
			Delete(h.getUsers(req.Info.Users...)); result.Error != nil {
			return result.Error
		}
	}
	return h.seeUpdate(repo, req, res)
}

func (h *Group) getUsers(in ...string) []model.User {
	list := make([]model.User, 0, len(in))
	for _, v := range in {
		list = append(list, model.User{
			Model: meta.Model{
				ID: v,
			},
		})
	}
	return list
}

// Invite Group
func (h *Group) Invite(ctx context.Context, req *pb.GroupModify, res *pb.GroupInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if _, err := h.check(repo, req.Visitor.UUID, req.Info.ID); err != nil {
		return err
	}

	if result := repo.
		Model(&model.Group{
			Model: meta.Model{
				ID: req.Info.ID,
			},
		}).
		Association("Users").
		Append(h.getUsers(req.Info.Users...)); result.Error != nil {
		return result.Error
	}
	return h.seeUpdate(repo, req, res)
}

// === Other ===

// Join Group
func (h *Group) Join(ctx context.Context, req *pb.GroupModify, res *pb.GroupInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.
		Model(&model.Group{
			Model: meta.Model{
				ID: req.Info.ID,
			},
		}).
		Association("Users").
		Append(h.getUsers(req.Visitor.UUID)); result.Error != nil {
		return result.Error
	}
	return h.seeUpdate(repo, req, res)
}

// Quit Group
func (h *Group) Quit(ctx context.Context, req *pb.GroupModify, res *pb.GroupInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.
		Model(&model.Group{
			Model: meta.Model{
				ID: req.Info.ID,
			},
		}).
		Association("Users").
		Delete(h.getUsers(req.Visitor.UUID)); result.Error != nil {
		return result.Error
	}
	return h.seeUpdate(repo, req, res)
}
