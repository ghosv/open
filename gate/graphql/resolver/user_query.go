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

// User Query
func (r *Resolver) User(ctx context.Context, args struct {
	UUID *string
}) (*UserResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acReadUserInfo) {
		return nil, meta.ErrAccessDenied
	}

	uuid := token.Base.UUID
	if args.UUID != nil {
		uuid = *args.UUID
	}
	user, err := loader.LoadUser(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return &UserResolver{ctx, acl, user}, nil
}

// Users Query
func (r *Resolver) Users(ctx context.Context, args struct {
	Word string
	Page *int32
	Size *int32
}) (*UserListResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check(acReadUserInfo) {
		return nil, meta.ErrAccessDenied
	}

	page := int32(1)
	if args.Page != nil {
		page = *args.Page
	}
	size := int32(5)
	if args.Page != nil {
		size = *args.Size
	}
	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfUser.Search(ctx, &selfPb.SearchForm{
		Word: args.Word,
		Page: page,
		Size: size,
	})
	if err != nil {
		return nil, err
	}

	return &UserListResolver{ctx, acl, res.Total, res.List}, nil
}
