package resolver

import (
	"context"

	"github.com/ghosv/open/gate/utils"
	pb "github.com/ghosv/open/plat/services/self/proto"
)

// UserListResolver of self
type UserListResolver struct {
	ctx   context.Context
	acl   *utils.ACL
	total int32
	list  []*pb.UserInfo
}

// Total of UserList
func (r *UserListResolver) Total() int32 {
	return r.total
}

// List of UserList
func (r *UserListResolver) List() *[]*UserResolver {
	results := make([]*UserResolver, 0, len(r.list))
	for _, v := range r.list {
		results = append(results, &UserResolver{r.ctx, r.acl, v})
	}
	return &results
}

// AppListResolver of self
type AppListResolver struct {
	ctx   context.Context
	acl   *utils.ACL
	total int32
	list  []*pb.AppInfo
}

// Total of AppList
func (r *AppListResolver) Total() int32 {
	return r.total
}

// List of AppList
func (r *AppListResolver) List() *[]*AppResolver {
	results := make([]*AppResolver, 0, len(r.list))
	for _, v := range r.list {
		results = append(results, &AppResolver{r.ctx, r.acl, v})
	}
	return &results
}

// OrgListResolver of self
type OrgListResolver struct {
	ctx   context.Context
	acl   *utils.ACL
	total int32
	list  []*pb.OrgInfo
}

// Total of OrgList
func (r *OrgListResolver) Total() int32 {
	return r.total
}

// List of OrgList
func (r *OrgListResolver) List() *[]*OrgResolver {
	results := make([]*OrgResolver, 0, len(r.list))
	for _, v := range r.list {
		results = append(results, &OrgResolver{r.ctx, r.acl, v})
	}
	return &results
}

// GroupListResolver of self
type GroupListResolver struct {
	ctx   context.Context
	acl   *utils.ACL
	total int32
	list  []*pb.GroupInfo
}

// Total of GroupList
func (r *GroupListResolver) Total() int32 {
	return r.total
}

// List of GroupList
func (r *GroupListResolver) List() *[]*GroupResolver {
	results := make([]*GroupResolver, 0, len(r.list))
	for _, v := range r.list {
		results = append(results, &GroupResolver{r.ctx, r.acl, v})
	}
	return &results
}
