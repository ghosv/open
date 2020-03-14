package resolver

import (
	"context"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/graphql/loader"
	"github.com/ghosv/open/gate/utils"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/core/proto"
	selfPb "github.com/ghosv/open/plat/services/self/proto"
)

// MyMarkApps Query
func (r *Resolver) MyMarkApps(ctx context.Context) ([]*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acReadUserAppList) {
		return nil, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.CoreApp.List(ctx, &pb.Identity{
		UUID: token.Base.UUID,
	})
	if err != nil {
		return nil, err
	}

	apps := make([]*AppResolver, 0, len(res.Data))
	for _, v := range res.Data {
		apps = append(apps, &AppResolver{ctx, acl, &selfPb.AppInfo{
			ID:      v.ID,
			Name:    v.Name,
			Icon:    v.Icon,
			Intro:   v.Intro,
			URL:     v.URL,
			OwnerID: v.OwnerID,
		}})
	}
	return apps, nil
}

// App Query
func (r *Resolver) App(ctx context.Context, args struct {
	ID string
}) (*AppResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acReadAppInfo) {
		return nil, meta.ErrAccessDenied
	}

	app, err := loader.LoadApp(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &AppResolver{ctx, acl, app}, nil
}

// Apps Query
func (r *Resolver) Apps(ctx context.Context, args struct {
	Word string
	Page *int
	Size *int
}) (*AppListResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acReadAppInfo) {
		return nil, meta.ErrAccessDenied
	}

	page := 1
	if args.Page != nil {
		page = *args.Page
	}
	size := 5
	if args.Page != nil {
		size = *args.Size
	}
	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfApp.Search(ctx, &selfPb.SearchForm{
		Word: args.Word,
		Page: int32(page),
		Size: int32(size),
	})
	if err != nil {
		return nil, err
	}

	return &AppListResolver{ctx, acl, res.Total, res.List}, nil
}
