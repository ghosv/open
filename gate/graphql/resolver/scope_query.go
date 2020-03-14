package resolver

import (
	"context"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/graphql/loader"
	"github.com/ghosv/open/gate/utils"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/core/proto"
)

// Scopes Query for 3rd
func (r *Resolver) Scopes(ctx context.Context, args struct {
	AppID     string
	AppKey    string
	AppSecret string
}) (*TokenScopesResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.CoreAuth.CheckTokenScope(ctx, &pb.AuthApp{
		AppID:     args.AppID,
		AppKey:    args.AppKey,
		AppSecret: args.AppSecret,
		Token:     token,
	})
	if err != nil {
		return nil, err
	}

	acl := utils.NewACL(token, meta.SrvSelf)
	return &TokenScopesResolver{ctx, acl, res}, nil
}

// TokenScopesResolver of core
type TokenScopesResolver struct {
	ctx     context.Context
	acl     *utils.ACL
	payload *pb.AuthScope
}

// User of Token
func (r *TokenScopesResolver) User() *UserResolver {
	if !r.acl.Check(acReadUserInfo) {
		return nil
	}
	user, err := loader.LoadUser(r.ctx, r.payload.UUID)
	if err != nil {
		return nil
	}
	return &UserResolver{r.ctx, r.acl, user}
}

// App of Token
func (r *TokenScopesResolver) App() *AppResolver {
	if !r.acl.Check(acReadAppInfo) {
		return nil
	}
	app, err := loader.LoadApp(r.ctx, r.payload.UUID)
	if err != nil {
		return nil
	}
	return &AppResolver{r.ctx, r.acl, app}
}

// Scopes of Token
func (r *TokenScopesResolver) Scopes() *[]string {
	return &r.payload.Scopes
}
