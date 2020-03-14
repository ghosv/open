package resolver

import (
	"context"

	"github.com/ghosv/open/gate/graphql/loader"
	"github.com/ghosv/open/gate/utils"
	pb "github.com/ghosv/open/plat/services/self/proto"
)

// GroupResolver of core
type GroupResolver struct {
	ctx   context.Context
	acl   *utils.ACL
	group *pb.GroupInfo
}

// ID of Group
func (r *GroupResolver) ID() string {
	return r.group.ID
}

// Name of Group
func (r *GroupResolver) Name() string {
	return r.group.Name
}

// Icon of Group
func (r *GroupResolver) Icon() string {
	return r.group.Icon
}

// Detail of Group
func (r *GroupResolver) Detail() string {
	return r.group.Detail
}

// Master of Group
func (r *GroupResolver) Master() *UserResolver {
	if r.group.MasterID == "" || !r.acl.Check(acReadUserInfo) {
		return nil
	}
	user, err := loader.LoadUser(r.ctx, r.group.MasterID)
	if err != nil {
		return nil
	}
	return &UserResolver{r.ctx, r.acl, user}
}

// Users of Group
func (r *GroupResolver) Users() *[]*UserResolver {
	if r.group.Users == nil || len(r.group.Users) == 0 ||
		!r.acl.Check(acReadUserInfo) {
		return nil
	}
	users, err := loader.LoadUsers(r.ctx, r.group.Users)
	if err != nil {
		return nil
	}

	results := make([]*UserResolver, 0, len(users))
	for _, v := range users {
		results = append(results, &UserResolver{r.ctx, r.acl, v})
	}
	return &results
}
