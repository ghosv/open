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

// CreateApp Mutation
func (r *Resolver) CreateApp(ctx context.Context, args struct {
	Name  string
	Icon  string
	Intro string
	URL   string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Create(ctx, &selfPb.AppModify{
		Visitor: acl.Visitor(),
		Name:    args.Name,
		Icon:    args.Icon,
		Intro:   args.Intro,
		URL:     args.URL,
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &AppResolver{ctx, acl, res}, nil
}

// DeleteApp Mutation
func (r *Resolver) DeleteApp(ctx context.Context, args struct {
	ID string
}) (bool, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return false, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	_, err := service.SelfApp.Delete(ctx, &selfPb.AppModify{
		Visitor: acl.Visitor(),
		ID:      args.ID,
	})
	if err != nil {
		return false, utils.MicroError(err)
	}
	return true, nil
}

// === Update ===

// ResetAppSecret Mutation
func (r *Resolver) ResetAppSecret(ctx context.Context, args struct {
	ID string
}) (*string, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Update(ctx, &selfPb.AppModify{
		Visitor:     acl.Visitor(),
		ID:          args.ID,
		ResetSecret: true,
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &res.Secret, nil
}

// UpdateApp Mutation
func (r *Resolver) UpdateApp(ctx context.Context, args struct {
	ID    string
	Name  *string
	Icon  *string
	Intro *string
	URL   *string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	var Name, Icon, Intro, URL string
	if args.Name != nil {
		Name = *args.Name
	}
	if args.Icon != nil {
		Icon = *args.Icon
	}
	if args.Intro != nil {
		Intro = *args.Intro
	}
	if args.URL != nil {
		URL = *args.URL
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Update(ctx, &selfPb.AppModify{
		Visitor: acl.Visitor(),
		ID:      args.ID,
		Name:    Name,
		Icon:    Icon,
		Intro:   Intro,
		URL:     URL,
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &AppResolver{ctx, acl, res}, nil
}

// AddAppScope Mutation
func (r *Resolver) AddAppScope(ctx context.Context, args struct {
	ID     string
	Name   string
	Detail string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Update(ctx, &selfPb.AppModify{
		Visitor:         acl.Visitor(),
		ID:              args.ID,
		ModifyScopeType: selfPb.ModifyScopeType_AddScope,
		Scope: &selfPb.AccessScope{
			Name:   args.Name,
			Detail: args.Detail,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &AppResolver{ctx, acl, res}, nil
}

// DelAppScope Mutation
func (r *Resolver) DelAppScope(ctx context.Context, args struct {
	ID   string
	Name string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Update(ctx, &selfPb.AppModify{
		Visitor:         acl.Visitor(),
		ID:              args.ID,
		ModifyScopeType: selfPb.ModifyScopeType_DelScope,
		Scope: &selfPb.AccessScope{
			Name: args.Name,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &AppResolver{ctx, acl, res}, nil
}

// AddAppAccess Mutation
func (r *Resolver) AddAppAccess(ctx context.Context, args struct {
	ID   string
	Name string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Update(ctx, &selfPb.AppModify{
		Visitor:         acl.Visitor(),
		ID:              args.ID,
		ModifyScopeType: selfPb.ModifyScopeType_AddAccess,
		Scope: &selfPb.AccessScope{
			Name: args.Name,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &AppResolver{ctx, acl, res}, nil
}

// DelAppAccess Mutation
func (r *Resolver) DelAppAccess(ctx context.Context, args struct {
	ID   string
	Name string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Update(ctx, &selfPb.AppModify{
		Visitor:         acl.Visitor(),
		ID:              args.ID,
		ModifyScopeType: selfPb.ModifyScopeType_DelAccess,
		Scope: &selfPb.AccessScope{
			Name: args.Name,
		},
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &AppResolver{ctx, acl, res}, nil
}

// === Other ===

// MarkApp Mutation
func (r *Resolver) MarkApp(ctx context.Context, args struct {
	ID string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Mark(ctx, &selfPb.AppModify{
		Visitor: acl.Visitor(),
		ID:      args.ID,
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &AppResolver{ctx, acl, res}, nil
}

// UnmarkApp Mutation
func (r *Resolver) UnmarkApp(ctx context.Context, args struct {
	ID string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Unmark(ctx, &selfPb.AppModify{
		Visitor: acl.Visitor(),
		ID:      args.ID,
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &AppResolver{ctx, acl, res}, nil
}
