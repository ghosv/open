package resolver

import (
	"context"

	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/utils"
	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/core/proto"
	notifyPb "github.com/ghosv/open/plat/services/notify/proto"
	selfPb "github.com/ghosv/open/plat/services/self/proto"
)

// UpdateUser Mutation
func (r *Resolver) UpdateUser(ctx context.Context, args struct {
	Nick     *string
	Avatar   *string
	Motto    *string
	Homepage *string
}) (*UserResolver, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	UUID := token.Base.UUID
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return nil, meta.ErrAccessDenied
	}

	var Nick, Avatar, Motto, Homepage string
	if args.Nick != nil {
		Nick = *args.Nick
	}
	if args.Avatar != nil {
		Avatar = *args.Avatar
	}
	if args.Motto != nil {
		Motto = *args.Motto
	}
	if args.Homepage != nil {
		Homepage = *args.Homepage
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	res, err := service.SelfUser.Update(ctx, &selfPb.UserModify{
		UUID:     UUID,
		Nick:     Nick,
		Avatar:   Avatar,
		Motto:    Motto,
		Homepage: Homepage,
	})
	if err != nil {
		return nil, utils.MicroError(err)
	}
	return &UserResolver{ctx, acl, res}, nil
}

// UpdateUserPass Mutation
func (r *Resolver) UpdateUserPass(ctx context.Context, args struct {
	Old  string
	Pass string
}) (bool, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	UUID := token.Base.UUID
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return false, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	_, err := service.SelfUser.Update(ctx, &selfPb.UserModify{
		UUID: UUID,
		Old:  args.Old,
		Pass: args.Pass,
	})
	if err != nil {
		return false, utils.MicroError(err)
	}
	return true, nil
}

// PostAuthCode Mutation
func (r *Resolver) PostAuthCode(ctx context.Context, args struct {
	Type string
	To   string
}) (bool, error) {
	var t notifyPb.PostType
	switch args.Type {
	case "Phone":
		t = notifyPb.PostType_CodePhone
	case "Email":
		t = notifyPb.PostType_CodeEmail
	case "GP":
		t = notifyPb.PostType_CodeGP
	default:
		return false, meta.ErrUnsupportedVerifyType
	}
	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	err := service.NotifyPostCode.Publish(ctx, &notifyPb.PostCode{
		Type: t,
		To:   args.To,
	})
	if err != nil {
		return false, utils.MicroError(err)
	}
	return true, nil
}

// UpdateUserBinding Mutation
func (r *Resolver) UpdateUserBinding(ctx context.Context, args struct {
	Type string
	To   string
	Code string
}) (bool, error) {
	token := ctx.Value(meta.KeyTokenPayload).(*pb.TokenPayload)
	UUID := token.Base.UUID
	acl := utils.NewACL(token, meta.SrvSelf)
	if !acl.Check() {
		return false, meta.ErrAccessDenied
	}

	service := ctx.Value(meta.KeyService).(*client.MicroClient)
	var err error
	switch args.Type {
	case "Phone":
		_, err = service.SelfUser.Update(ctx, &selfPb.UserModify{
			UUID:      UUID,
			Phone:     args.To,
			PhoneCode: args.Code,
		})
	case "Email":
		_, err = service.SelfUser.Update(ctx, &selfPb.UserModify{
			UUID:      UUID,
			Email:     args.To,
			EmailCode: args.Code,
		})
	default:
		return false, meta.ErrUnsupportedVerifyType
	}

	if err != nil {
		return false, utils.MicroError(err)
	}
	return true, nil
}
