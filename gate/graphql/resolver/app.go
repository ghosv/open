package resolver

import (
	"context"

	"github.com/ghosv/open/gate/graphql/loader"
	"github.com/ghosv/open/gate/utils"
	pb "github.com/ghosv/open/plat/services/self/proto"
)

// === public ===

// AppResolver of self
type AppResolver struct {
	ctx context.Context
	acl *utils.ACL
	app *pb.AppInfo
}

// ID of App
func (r *AppResolver) ID() string {
	return r.app.ID
}

// Name of App
func (r *AppResolver) Name() string {
	return r.app.Name
}

// Icon of App
func (r *AppResolver) Icon() string {
	return r.app.Icon
}

// Intro of App
func (r *AppResolver) Intro() string {
	return r.app.Intro
}

// URL of App
func (r *AppResolver) URL() string {
	return r.app.URL
}

// Owner of App
func (r *AppResolver) Owner() *UserResolver {
	if !r.acl.Check(acReadUserInfo) {
		return nil
	}
	user, err := loader.LoadUser(r.ctx, r.app.OwnerID)
	if err != nil {
		return nil
	}
	return &UserResolver{r.ctx, r.acl, user}
}

// Scopes of App
func (r *AppResolver) Scopes() *[]*AccessScopeResolver {
	return r.scopes(r.app.Scopes)
}

func (r *AppResolver) scopes(scopes []*pb.AccessScope) *[]*AccessScopeResolver {
	if scopes == nil {
		return nil
	}
	list := make([]*AccessScopeResolver, 0, len(scopes))
	for _, v := range scopes {
		list = append(list, &AccessScopeResolver{r.ctx, r.acl, v})
	}
	return &list
}

// AccessList of App
func (r *AppResolver) AccessList() *[]*AccessScopeResolver {
	return r.scopes(r.app.AccessList)
}

// === private ===

// Key of App
func (r *AppResolver) Key() *string {
	if !r.acl.Display(r.app.OwnerID) {
		return nil
	}
	return &r.app.Key
}

// Secret of App
func (r *AppResolver) Secret() *string {
	if !r.acl.Display(r.app.OwnerID) {
		return nil
	}
	return &r.app.Secret
}

func (r *AppResolver) users(keys []string) *[]*UserResolver {
	if keys == nil || !r.acl.Check(acReadUserInfo) || !r.acl.Display(r.app.OwnerID) {
		return nil
	}
	users, err := loader.LoadUsers(r.ctx, keys)
	if err != nil {
		return nil
	}
	results := make([]*UserResolver, 0, len(users))
	for _, v := range users {
		results = append(results, &UserResolver{r.ctx, r.acl, v})
	}
	return &results
}

// Managers of App
func (r *AppResolver) Managers() *[]*UserResolver {
	return r.users(r.app.Managers)
}

// Developers of App
func (r *AppResolver) Developers() *[]*UserResolver {
	return r.users(r.app.Developers)
}

// Users of App
func (r *AppResolver) Users() *[]*UserResolver {
	return r.users(r.app.Users)
}
