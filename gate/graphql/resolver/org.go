package resolver

import (
	"context"

	"github.com/ghosv/open/gate/graphql/loader"
	"github.com/ghosv/open/gate/utils"
	pb "github.com/ghosv/open/plat/services/self/proto"
)

// OrgResolver of core
type OrgResolver struct {
	ctx context.Context
	acl *utils.ACL
	org *pb.OrgInfo
}

// ID of Org
func (r *OrgResolver) ID() string {
	return r.org.ID
}

// Name of Org
func (r *OrgResolver) Name() string {
	return r.org.Name
}

// Icon of Org
func (r *OrgResolver) Icon() string {
	return r.org.Icon
}

// Detail of Org
func (r *OrgResolver) Detail() string {
	return r.org.Detail
}

// Master of Org
func (r *OrgResolver) Master() *UserResolver {
	if r.org.MasterID == "" || !r.acl.Check(acReadUserInfo) {
		return nil
	}
	user, err := loader.LoadUser(r.ctx, r.org.MasterID)
	if err != nil {
		return nil
	}
	return &UserResolver{r.ctx, r.acl, user}
}

// Users of Org
func (r *OrgResolver) Users() *[]*UserResolver {
	if r.org.Users == nil || len(r.org.Users) == 0 ||
		!r.acl.Check(acReadUserInfo) {
		return nil
	}
	users, err := loader.LoadUsers(r.ctx, r.org.Users)
	if err != nil {
		return nil
	}

	results := make([]*UserResolver, 0, len(users))
	for _, v := range users {
		results = append(results, &UserResolver{r.ctx, r.acl, v})
	}
	return &results
}

// Father of Org
func (r *OrgResolver) Father() *OrgResolver {
	if r.org.FatherID == "" {
		return nil
	}
	org, err := loader.LoadOrg(r.ctx, r.org.FatherID)
	if err != nil {
		return nil
	}
	return &OrgResolver{r.ctx, r.acl, org}
}

// Children of Org
func (r *OrgResolver) Children() *[]*OrgResolver {
	if r.org.ChildrenID == nil || len(r.org.ChildrenID) == 0 {
		return nil
	}
	orgs, err := loader.LoadOrgs(r.ctx, r.org.ChildrenID)
	if err != nil {
		return nil
	}

	results := make([]*OrgResolver, 0, len(orgs))
	for _, v := range orgs {
		results = append(results, &OrgResolver{r.ctx, r.acl, v})
	}
	return &results
}
