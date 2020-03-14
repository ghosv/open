package handler

import (
	"context"

	"github.com/ghosv/open/meta"
	pb "github.com/ghosv/open/plat/services/notify/proto"
	"github.com/ghosv/open/plat/services/notify/subscriber"
	"github.com/go-redis/redis/v7"
)

// Verify Handler
type Verify struct{}

// CodeMatch for core
func (h *Verify) CodeMatch(ctx context.Context, req *pb.VerifyCodeGroup, res *pb.VerifyCodeGroup) error {
	// TODO: impl
	client := ctx.Value(meta.KeyRepoRedis).(*redis.Client)
	res.Codes = req.Codes
	for _, v := range req.Codes {
		k := subscriber.KeyGroupRedisVerifyCodes[v.Type] + v.To
		code, e := client.Get(k).Result()
		v.Match = e == nil && v.Code == code
		if v.Match {
			client.Del(k)
		}
	}
	return nil
}

// CodeDebug <tmp>
// TODO: remove
func (h *Verify) CodeDebug(ctx context.Context, req *pb.VerifyCodeGroup, res *pb.VerifyCodeGroup) error {
	client := ctx.Value(meta.KeyRepoRedis).(*redis.Client)
	res.Codes = req.Codes
	for _, v := range req.Codes {
		k := subscriber.KeyGroupRedisVerifyCodes[v.Type] + v.To
		v.Code, _ = client.Get(k).Result()
	}
	return nil
}
