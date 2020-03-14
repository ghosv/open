package handler

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/ghosv/open/meta"
	"github.com/ghosv/open/plat/services/core/initialize"
	model "github.com/ghosv/open/plat/services/core/model"
	pb "github.com/ghosv/open/plat/services/core/proto"
	notifyPb "github.com/ghosv/open/plat/services/notify/proto"
	"github.com/ghosv/open/plat/utils"
)

// User Handler
type User struct{}

// Login Action
func (h *User) Login(ctx context.Context, req *pb.Credential, res *pb.Token) error {
	name, phone, email, ok := utils.CheckUsername(req.NamePhoneEmail)
	if !ok {
		return meta.ErrWrongUsername
	}
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	u := model.User{
		Name:  name,
		Phone: phone,
		Email: email,
		Pass:  req.Pass,
	}
	if result := repo.Where(&u).First(&u); result.Error != nil {
		return meta.ErrCredential
	}

	token := NewToken()
	token.Payload = u.TokenPayload()
	var e error
	if res.Str, e = token.Sign(initialize.DefaultJWTTerm); e != nil {
		return e
	}
	return nil
}

// Register Action
func (h *User) Register(ctx context.Context, req *pb.Identity, res *pb.Token) error {
	name, _, _, ok := utils.CheckUsername(req.Name)
	if !(ok && name == req.Name) ||
		!utils.IsPhone(req.Phone) ||
		!utils.IsEmail(req.Email) {
		return meta.ErrWrongUsername
	}

	// verify phone & email
	c := ctx.Value(meta.KeyClientNotifyVerify).(notifyPb.VerifyService)
	res0, err := c.CodeMatch(ctx, &notifyPb.VerifyCodeGroup{
		Codes: []*notifyPb.VerifyCode{
			&notifyPb.VerifyCode{
				Type: notifyPb.PostType_CodePhone,
				To:   req.Phone,
				Code: req.PhoneCode,
			},
			&notifyPb.VerifyCode{
				Type: notifyPb.PostType_CodeEmail,
				To:   req.Email,
				Code: req.EmailCode,
			},
		},
	})
	if err != nil {
		return meta.ErrInvalidCode
	}
	for _, v := range res0.Codes {
		if !v.Match {
			return meta.ErrInvalidCode
		}
	}

	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	u := &model.User{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
		Pass:  req.Pass,

		Nick: req.Name,
	}
	if result := repo.Create(u); result.Error != nil {
		return meta.ErrUserHasExist
	}

	token := NewToken()
	token.Payload = u.TokenPayload()
	var e error
	if res.Str, e = token.Sign(initialize.DefaultJWTTerm); e != nil {
		return e
	}
	return nil
}

// Check Token
func (h *User) Check(ctx context.Context, req *pb.Token, res *pb.Token) error {
	token, err := ValidToken(req.Str)
	if err != nil {
		return err
	}
	t, err := FreshToken(token,
		initialize.DefaultJWTFresh, initialize.DefaultJWTTerm)
	res.Str = t
	res.Payload = token.Payload
	return nil
}
