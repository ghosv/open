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

// CreateOrg Mutation
func (r *Resolver) CreateOrg(ctx context.Context, args struct {
	Father *string
	Name   string
	Icon   string
	Detail string
	Master *string
}) (*OrgResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	Father := ""
	if args.Father != nil {
		Father = *args.Father
	}
	Master := token.Base.UUID
	if args.Master != nil {
		Master = *args.Master
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfOrg.Create(ctx, &selfPb.OrgModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.OrgInfo{
			FatherID: Father,
			Name:     args.Name,
			Icon:     args.Icon,
			Detail:   args.Detail,
			MasterID: Master,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &OrgResolver{ctx, acl, res}, nil
}

// DeleteOrg Mutation
func (r *Resolver) DeleteOrg(ctx context.Context, args struct {
	ID string
}) (bool, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return false, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	_, err := service.SelfOrg.Delete(ctx, &selfPb.OrgModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.OrgInfo{
			ID: args.ID,
		},
	})
	if err != nil {
		return false, utils.MicroError(err)
	}
	return true, nil
}

// === Update ===

// UpdateOrg Mutation
func (r *Resolver) UpdateOrg(ctx context.Context, args struct {
	ID       string
	Name     *string
	Icon     *string
	Detail   *string
	Master   *string
	DelUsers *[]string
}) (*OrgResolver, error) {
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
	res, err := service.SelfOrg.Update(ctx, &selfPb.OrgModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.OrgInfo{
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
	return &OrgResolver{ctx, acl, res}, nil
}

// InviteJoinOrg Mutation
func (r *Resolver) InviteJoinOrg(ctx context.Context, args struct {
	ID    string
	Users []string
}) (*OrgResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfOrg.Invite(ctx, &selfPb.OrgModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.OrgInfo{
			ID:    args.ID,
			Users: args.Users,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &OrgResolver{ctx, acl, res}, nil
}

// === Other ===

// JoinOrg Mutation
func (r *Resolver) JoinOrg(ctx context.Context, args struct {
	ID string
}) (*OrgResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfOrg.Join(ctx, &selfPb.OrgModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.OrgInfo{
			ID: args.ID,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &OrgResolver{ctx, acl, res}, nil
}

// QuitOrg Mutation
func (r *Resolver) QuitOrg(ctx context.Context, args struct {
	ID string
}) (*OrgResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acWriteGroupInfo) {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfOrg.Quit(ctx, &selfPb.OrgModify{
		Visitor: acl.Visitor(),
		Info: &selfPb.OrgInfo{
			ID: args.ID,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &OrgResolver{ctx, acl, res}, nil
}
