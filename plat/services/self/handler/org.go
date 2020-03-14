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

// Org Handler
type Org struct{}

// === Query ===

// BatchFind Org
func (h *Org) BatchFind(ctx context.Context, req *pb.BatchID, res *pb.OrgMap) error {
	var orgs []model.Org
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.Where(req.UUID).Find(&orgs); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	data := make(map[string]*pb.OrgInfo, len(orgs))
	for _, v := range orgs {
		d, e := h.loadData(repo, &v)
		if e != nil {
			return e
		}
		data[v.ID] = d
	}
	res.Data = data
	return nil
}

func (h *Org) loadData(repo *gorm.DB, v *model.Org) (*pb.OrgInfo, error) {
	var childrenID, users []string
	if result := repo.
		Table("orgs").
		Where("father_id = ?", v.ID).
		Pluck("id", &childrenID); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}
	if result := repo.
		Table("org_users").
		Where("org_id = ?", v.ID).
		Pluck("user_id", &users); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}
	return &pb.OrgInfo{
		ID:         v.ID,
		FatherID:   v.FatherID,
		ChildrenID: childrenID,
		Name:       v.Name,
		Icon:       v.Icon,
		Detail:     v.Detail,
		MasterID:   v.MasterID,
		Users:      users,
	}, nil
}

// Search Org
func (h *Org) Search(ctx context.Context, req *pb.SearchForm, res *pb.OrgList) error {
	word, page, size := req.Word, req.Page, req.Size
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	s := repo.
		Model(&model.Org{}).
		Where("name LIKE ?", "%"+word+"%")
	if result := s.Count(&res.Total); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	var orgs []model.Org
	if result := s.Limit(size).Offset((page - 1) * size).
		Find(&orgs); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	res.List = make([]*pb.OrgInfo, 0, len(orgs))
	for _, v := range orgs {
		d, e := h.loadData(repo, &v)
		if e != nil {
			return e
		}
		res.List = append(res.List, d)
	}
	return nil
}

// === Mutation ===

func (h *Org) check(repo *gorm.DB, visitor, id string) (*model.Org, error) {
	if id == "" {
		// return nil, meta.ErrOrgNotExist
		return nil, meta.ErrAccessDenied
	}
	data := model.Org{
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
	if _, err := h.check(repo, visitor, data.FatherID); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *Org) seeUpdate(repo *gorm.DB, req *pb.OrgModify, res *pb.OrgInfo) error {
	data := model.Org{
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

// Create Org
func (h *Org) Create(ctx context.Context, req *pb.OrgModify, res *pb.OrgInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if req.Info.FatherID != "" {
		if _, err := h.check(repo, req.Visitor.UUID,
			req.Info.FatherID); err != nil {
			return err
		}
	}

	data := &model.Org{
		FatherID: req.Info.FatherID,
		Name:     req.Info.Name,
		Icon:     req.Info.Icon,
		Detail:   req.Info.Detail,
		MasterID: req.Info.MasterID,
	}
	if result := repo.Create(data); result.Error != nil {
		return meta.ErrOrgHasExist
	}
	res0, err := h.loadData(repo, data)
	if err != nil {
		return err
	}
	*res = *res0
	return nil
}

// Delete Org
func (h *Org) Delete(ctx context.Context, req *pb.OrgModify, res *metaPb.None) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if _, err := h.check(repo, req.Visitor.UUID, req.Info.ID); err != nil {
		return err
	}

	if result := repo.Delete(&model.Org{
		Model: meta.Model{
			ID: req.Info.ID,
		},
	}); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	return nil
}

// Update Org
func (h *Org) Update(ctx context.Context, req *pb.OrgModify, res *pb.OrgInfo) error {
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

func (h *Org) getUsers(in ...string) []model.User {
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

// Invite Org
func (h *Org) Invite(ctx context.Context, req *pb.OrgModify, res *pb.OrgInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if _, err := h.check(repo, req.Visitor.UUID, req.Info.ID); err != nil {
		return err
	}

	if result := repo.
		Model(&model.Org{
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

// Join Org
func (h *Org) Join(ctx context.Context, req *pb.OrgModify, res *pb.OrgInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.
		Model(&model.Org{
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

// Quit Org
func (h *Org) Quit(ctx context.Context, req *pb.OrgModify, res *pb.OrgInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.
		Model(&model.Org{
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
