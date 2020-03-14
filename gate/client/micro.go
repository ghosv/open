package client

import (
	"log"

	"github.com/ghosv/open/gate/conf"
	"github.com/ghosv/open/meta"

	corePb "github.com/ghosv/open/plat/services/core/proto"
	notifyPb "github.com/ghosv/open/plat/services/notify/proto"
	selfPb "github.com/ghosv/open/plat/services/self/proto"

	"github.com/micro/go-micro/v2"
)

// Default Services
var Default = new(MicroClient)

// MicroClient for gate
type MicroClient struct {
	NotifyVerify   notifyPb.VerifyService
	NotifyPostCode micro.Publisher

	CoreUser corePb.UserService
	CoreApp  corePb.AppService
	CoreAuth corePb.AuthService

	SelfUser  selfPb.UserService
	SelfApp   selfPb.AppService
	SelfOrg   selfPb.OrgService
	SelfGroup selfPb.GroupService
}

// Init Client
func (c *MicroClient) Init(service micro.Service) {
	c.NotifyVerify = notifyPb.NewVerifyService(meta.SrvNotify, service.Client())
	c.NotifyPostCode = micro.NewPublisher(meta.TopicNotifyPostCode, service.Client())

	c.CoreUser = corePb.NewUserService(meta.SrvCore, service.Client())
	c.CoreApp = corePb.NewAppService(meta.SrvCore, service.Client())
	c.CoreAuth = corePb.NewAuthService(meta.SrvCore, service.Client())

	c.SelfUser = selfPb.NewUserService(meta.SrvSelf, service.Client())
	c.SelfApp = selfPb.NewAppService(meta.SrvSelf, service.Client())
	c.SelfOrg = selfPb.NewOrgService(meta.SrvSelf, service.Client())
	c.SelfGroup = selfPb.NewGroupService(meta.SrvSelf, service.Client())

	log.Println(conf.ReadyMicroClient)
}
