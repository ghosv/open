package handler

import (
	"context"

	"github.com/ghosv/open/meta"
	"github.com/ghosv/open/plat/services/core/model"
	notifyPb "github.com/ghosv/open/plat/services/notify/proto"
	pb "github.com/ghosv/open/plat/services/self/proto"
	"github.com/ghosv/open/plat/utils"
	"github.com/jinzhu/gorm"
)

// User Handler
type User struct{}

// === Query ===

// BatchFind User
func (h *User) BatchFind(ctx context.Context, req *pb.BatchID, res *pb.UserMap) error {
	var users []model.User
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.Where(req.UUID).Find(&users); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	data := make(map[string]*pb.UserInfo, len(users))
	for _, v := range users {
		d, e := h.loadData(repo, &v)
		if e != nil {
			return e
		}
		data[v.ID] = d
	}
	res.Data = data
	return nil
}

func (h *User) loadData(repo *gorm.DB, v *model.User) (*pb.UserInfo, error) {
	var myApps, apps []string
	if result := repo.
		Table("apps").
		Where(&model.App{OwnerID: v.ID}).
		Pluck("id", &myApps); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}
	if result := repo.
		Table("app_users").
		Where("user_id = ?", v.ID).
		Pluck("app_id", &apps); result.Error != nil {
		return nil, utils.RepoErrorFilter(result)
	}
	return &pb.UserInfo{
		UUID:     v.ID,
		Nick:     v.Nick,
		Avatar:   v.Avatar,
		Motto:    v.Motto,
		Homepage: v.Homepage,
		MyApps:   myApps,

		Name:  v.Name,
		Phone: v.Phone,
		Email: v.Email,
		Apps:  apps,
	}, nil
}

// Search User
func (h *User) Search(ctx context.Context, req *pb.SearchForm, res *pb.UserList) error {
	word, page, size := req.Word, req.Page, req.Size
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	s := repo.
		Model(&model.User{}).
		Where("nick LIKE ?", "%"+word+"%") // Or("name LIKE ?", "%"+word+"%").
	if result := s.Count(&res.Total); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	var users []model.User
	if result := s.Limit(size).Offset((page - 1) * size).
		Find(&users); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	res.List = make([]*pb.UserInfo, 0, len(users))
	for _, v := range users {
		d, e := h.loadData(repo, &v)
		if e != nil {
			return e
		}
		res.List = append(res.List, d)
	}
	return nil
}

// === Mutation ===

// Update User
func (h *User) Update(ctx context.Context, req *pb.UserModify, res *pb.UserInfo) error {
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	u := model.User{
		Model: meta.Model{
			ID: req.UUID,
		},
	}
	// TODO: test
	if result := repo. /*Where(&u).*/ First(&u); result.Error != nil {
		return meta.ErrCredential
	}

	if req.Old != "" && req.Old != u.Pass {
		return meta.ErrCredential
	}

	codes := make([]*notifyPb.VerifyCode, 0, 2)
	if req.Phone != "" {
		codes = append(codes, &notifyPb.VerifyCode{
			Type: notifyPb.PostType_CodePhone,
			To:   req.Phone,
			Code: req.PhoneCode,
		})
	}
	if req.Email != "" {
		codes = append(codes, &notifyPb.VerifyCode{
			Type: notifyPb.PostType_CodeEmail,
			To:   req.Email,
			Code: req.EmailCode,
		})
	}
	if len(codes) > 0 {
		// verify phone & email
		c := ctx.Value(meta.KeyClientNotifyVerify).(notifyPb.VerifyService)
		res0, err := c.CodeMatch(ctx, &notifyPb.VerifyCodeGroup{
			Codes: codes,
		})
		if err != nil {
			return meta.ErrInvalidCode
		}
		for _, v := range res0.Codes {
			if !v.Match {
				return meta.ErrInvalidCode
			}
		}
	}

	u.Name = ""
	u.Pass = req.Pass
	u.Phone = req.Phone
	u.Email = req.Email

	u.Nick = req.Nick
	u.Avatar = req.Avatar
	u.Motto = req.Motto
	u.Homepage = req.Homepage

	if result := repo.Model(&u).Updates(&u); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}

	if result := repo.Where(&model.User{
		Model: meta.Model{
			ID: req.UUID,
		},
	}).First(&u); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	res0, err := h.loadData(repo, &u)
	if err != nil {
		return err
	}
	*res = *res0
	return nil
}
