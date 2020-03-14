package resolver

import (
	"context"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/utils"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/core/proto"
	selfPb "github.com/ghosv/open/plat/services/self/proto"
)

// === Create & Delete ===

// CreateGroup Mutation
func (r *Resolver) CreateGroup(ctx context.Context, args struct {
	Name   string
	Icon   string
	Detail string
}) (*GroupResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfGroup.Create(ctx, &selfPb.GroupModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.GroupInfo{
			Name:   args.Name,
			Icon:   args.Icon,
			Detail: args.Detail,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &GroupResolver{ctx, acl, res}, nil
}

// DeleteGroup Mutation
func (r *Resolver) DeleteGroup(ctx context.Context, args struct {
	ID string
}) (bool, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return false, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	_, err := service.SelfGroup.Delete(ctx, &selfPb.GroupModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.GroupInfo{
			ID: args.ID,
		},
	})
	if err != nil {
		return false, utils.MicroError(err)
	}
	return true, nil
}

// === Update ===

// UpdateGroup Mutation
func (r *Resolver) UpdateGroup(ctx context.Context, args struct {
	ID       string
	Name     *string
	Icon     *string
	Detail   *string
	Master   *string
	DelUsers *[]string
}) (*GroupResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	var Name, Icon, Detail, Master string
	if args.Name != nil {
		Name = *args.Name
	}
	if args.Icon != nil {
		Icon = *args.Icon
	}
	if args.Detail != nil {
		Detail = *args.Detail
	}
	if args.Master != nil {
		Master = *args.Master
	}
	DelUsers := []string{}
	if args.DelUsers != nil {
		DelUsers = *args.DelUsers
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfGroup.Update(ctx, &selfPb.GroupModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.GroupInfo{
			ID:       args.ID,
			Name:     Name,
			Icon:     Icon,
			Detail:   Detail,
			MasterID: Master,
			Users:    DelUsers,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &GroupResolver{ctx, acl, res}, nil
}

// InviteJoinGroup Mutation
func (r *Resolver) InviteJoinGroup(ctx context.Context, args struct {
	ID    string
	Users []string
}) (*GroupResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfGroup.Invite(ctx, &selfPb.GroupModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.GroupInfo{
			ID:    args.ID,
			Users: args.Users,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &GroupResolver{ctx, acl, res}, nil
}

// === Other ===

// JoinGroup Mutation
func (r *Resolver) JoinGroup(ctx context.Context, args struct {
	ID string
}) (*GroupResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfGroup.Join(ctx, &selfPb.GroupModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.GroupInfo{
			ID: args.ID,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &GroupResolver{ctx, acl, res}, nil
}

// QuitGroup Mutation
func (r *Resolver) QuitGroup(ctx context.Context, args struct {
	ID string
}) (*GroupResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfGroup.Quit(ctx, &selfPb.GroupModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.GroupInfo{
			ID: args.ID,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &GroupResolver{ctx, acl, res}, nil
}
