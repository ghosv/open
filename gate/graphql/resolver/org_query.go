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

// Org Query
func (r *Resolver) Org(ctx context.Context, args struct {
	ID string
}) (*OrgResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acReadOrgInfo) {
		return nil, meta.ErrAccessDenied
	}

	org, err := loader.LoadOrg(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &OrgResolver{ctx, acl, org}, nil
}

// Orgs Query
func (r *Resolver) Orgs(ctx context.Context, args struct {
	Word string
	Page *int
	Size *int
}) (*OrgListResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acReadOrgInfo) {
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
	res, err := service.SelfOrg.Search(ctx, &selfPb.SearchForm{
		Word: args.Word,
		Page: int32(page),
		Size: int32(size),
	})
	if err != nil {
		return nil, err
	}

	return &OrgListResolver{ctx, acl, res.Total, res.List}, nil
}
