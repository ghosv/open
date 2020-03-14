package resolver

import (
	"context"

	"github.com/ghosv/open/gate/graphql/loader"
	"github.com/ghosv/open/gate/utils"
	pb "github.com/ghosv/open/plat/services/self/proto"
)

// UserResolver of core
type UserResolver struct {
	ctx  context.Context
	acl  *utils.ACL
	user *pb.UserInfo
}

// === public ===

// UUID of User
func (r *UserResolver) UUID() string {
	return r.user.UUID
}

// Nick of User
func (r *UserResolver) Nick() string {
	return r.user.Nick
}

// Avatar of User
func (r *UserResolver) Avatar() string {
	return r.user.Avatar
}

// Motto of User
func (r *UserResolver) Motto() string {
	return r.user.Motto
}

// Homepage of User
func (r *UserResolver) Homepage() string {
	return r.user.Homepage
}

// MyApps of User
func (r *UserResolver) MyApps() *[]*AppResolver {
	if r.user.MyApps == nil || !r.acl.Check(acReadAppInfo) {
		return nil
	}
	return r.apps(r.user.MyApps)
}

func (r *UserResolver) apps(keys []string) *[]*AppResolver {
	apps, err := loader.LoadApps(r.ctx, keys)
	if err != nil {
		return nil
	}
	results := make([]*AppResolver, 0, len(apps))
	for _, v := range apps {
		results = append(results, &AppResolver{r.ctx, r.acl, v})
	}
	return &results
}

// === private ===

// Name of User
func (r *UserResolver) Name() *string {
	if !r.acl.Display(r.UUID()) {
		return nil
	}
	return &r.user.Name
}

// Phone of User
func (r *UserResolver) Phone() *string {
	if !r.acl.Display(r.UUID()) {
		return nil
	}
	return &r.user.Phone
}

// Email of User
func (r *UserResolver) Email() *string {
	if !r.acl.Display(r.UUID()) {
		return nil
	}
	return &r.user.Email
}

// Apps of User
func (r *UserResolver) Apps() *[]*AppResolver {
	if r.user.Apps == nil || !r.acl.Check(acReadAppInfo) || !r.acl.Display(r.UUID()) {
		return nil
	}
	return r.apps(r.user.Apps)
}
