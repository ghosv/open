package handler

import (
	"context"
	"strings"

	proto "github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"

	"github.com/ghosv/open/meta"
	metaPb "github.com/ghosv/open/meta/proto"
	"github.com/ghosv/open/plat/services/core/initialize"
	model "github.com/ghosv/open/plat/services/core/model"
	pb "github.com/ghosv/open/plat/services/core/proto"
	"github.com/ghosv/open/plat/utils"

	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
)

// Auth Handler
type Auth struct{}

// === Authorize ===

const redisKeyAuthCode = "auth.oauth2.code-"

// Authorize Scopes
func (h *Auth) Authorize(ctx context.Context, req *pb.AuthRequest, res *pb.AuthResponse) error {
	token, err := ValidToken(req.Token)
	if err != nil {
		return err
	}
	if req.Type != pb.AuthType_Code {
		return meta.ErrUnsupportedAuthType
	}

	app := &model.App{
		Model: meta.Model{ID: req.AppID},
	}
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.First(app); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	var scopes []model.AccessScope
	if result := repo.Model(app).Related(&scopes, "AccessList"); result.Error != nil {
		if !result.RecordNotFound() {
			return result.Error
		}
	}
	_scopes := make([]*pb.AccessScope, 0, len(scopes))
	for _, v := range scopes {
		_scopes = append(_scopes, &pb.AccessScope{
			Name:   v.Name,
			Detail: v.Detail,
		})
	}

	code := uuid.NewV4().String()

	client := ctx.Value(meta.KeyRepoRedis).(*redis.Client)
	client.Set(redisKeyAuthCode+code+".app",
		app.ID+"."+app.Key+"."+app.Secret,
		initialize.DefaultAuthCodeTime)
	client.Set(redisKeyAuthCode+code+".user",
		proto.MarshalTextString(token.Payload),
		initialize.DefaultAuthCodeTime)

	res.Code = code
	res.Info = &pb.AppInfo{
		ID:         app.ID,
		Name:       app.Name,
		Icon:       app.Icon,
		Intro:      app.Intro,
		URL:        app.URL,
		AccessList: _scopes,
	}
	return nil
}

// AuthorizeConfirm Scopes
func (h *Auth) AuthorizeConfirm(ctx context.Context, req *pb.AuthConfirm, res *metaPb.None) error {
	client := ctx.Value(meta.KeyRepoRedis).(*redis.Client)
	client.Set(redisKeyAuthCode+req.Code+".scopes",
		proto.MarshalTextString(req),
		initialize.DefaultAuthCodeTime)
	return nil
}

// AuthToken for App
func (h *Auth) AuthToken(ctx context.Context, req *pb.AuthApp, res *pb.Token) error {
	client := ctx.Value(meta.KeyRepoRedis).(*redis.Client)
	app, e := client.Get(redisKeyAuthCode + req.Code + ".app").Result()
	if e == redis.Nil {
		return meta.ErrInvalidCode
	}
	if req.AppID+"."+req.AppKey+"."+req.AppSecret != app {
		return meta.ErrAuthenticationFailed
	}

	_token, e := client.Get(redisKeyAuthCode + req.Code + ".user").Result()
	if e == redis.Nil {
		return meta.ErrInvalidCode
	}
	token := NewToken()
	token.Payload = new(pb.TokenPayload)
	proto.UnmarshalText(_token, token.Payload)

	_scopes, e := client.Get(redisKeyAuthCode + req.Code + ".scopes").Result()
	if e == redis.Nil {
		return meta.ErrInvalidCode
	}
	tmp := new(pb.AuthConfirm)
	proto.UnmarshalText(_scopes, tmp)

	token.Payload.AppID = req.AppID
	token.Payload.Scopes = tmp.Scopes
	res.Str, _ = token.Sign(initialize.DefaultJWTTerm)

	client.Del(redisKeyAuthCode + req.Code + ".app")
	client.Del(redisKeyAuthCode + req.Code + ".user")
	client.Del(redisKeyAuthCode + req.Code + ".scopes")
	return nil
}

// CheckTokenScope Scopes
func (h *Auth) CheckTokenScope(ctx context.Context, req *pb.AuthApp, res *pb.AuthScope) error {
	app := &model.App{
		Model: meta.Model{ID: req.AppID},
	}
	repo := ctx.Value(meta.KeyRepo).(*gorm.DB)
	if result := repo.First(app); result.Error != nil {
		return utils.RepoErrorFilter(result)
	}
	if app.Key != req.AppKey || app.Secret != req.AppSecret {
		return meta.ErrAuthenticationFailed
	}

	res.AppID = req.Token.AppID
	res.UUID = req.Token.Base.UUID
	res.Scopes = make([]string, 0, len(req.Token.Scopes))
	for _, v := range req.Token.Scopes {
		if !strings.HasPrefix(v.Name, req.AppID+":") {
			continue
		}
		res.Scopes = append(res.Scopes,
			strings.Replace(v.Name, req.AppID+":", "", 1))
	}
	return nil
}
