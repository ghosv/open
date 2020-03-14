package resolver

import (
	"context"

	"github.com/ghosv/open/gate/graphql/loader"
	"github.com/ghosv/open/gate/utils"
	pb "github.com/ghosv/open/plat/services/self/proto"
)

// AccessScopeResolver of self
type AccessScopeResolver struct {
	ctx   context.Context
	acl   *utils.ACL
	scope *pb.AccessScope
}

// App of AccessScope
func (r *AccessScopeResolver) App() *AppResolver {
	if !r.acl.Check(acReadAppInfo) {
		return nil
	}
	app, err := loader.LoadApp(r.ctx, r.scope.AppID)
	if err != nil {
		return nil
	}
	return &AppResolver{r.ctx, r.acl, app}
}

// Name of AccessScope
func (r *AccessScopeResolver) Name() string {
	return r.scope.Name
}

// Detail of AccessScope
func (r *AccessScopeResolver) Detail() string {
	return r.scope.Detail
}
