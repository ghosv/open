package subscriber

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/ghosv/open/meta"
	"github.com/ghosv/open/plat/services/notify/initialize"
	pb "github.com/ghosv/open/plat/services/notify/proto"
	"github.com/go-redis/redis/v7"
)

// KeyGroupRedisVerifyCodes for post
var KeyGroupRedisVerifyCodes = map[pb.PostType]string{
	pb.PostType_CodePhone: meta.KeyRedisPhoneCode,
	pb.PostType_CodeEmail: meta.KeyRedisEmailCode,
	pb.PostType_CodeGP:    meta.KeyRedisGhostIMCode,
}

func newCode() string {
	code := rand.Int() % 1_000_000
	s := strconv.Itoa(code)
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}

// PostCode Sub
type PostCode struct{}

// Handle PostCode
func (s *PostCode) Handle(ctx context.Context, msg *pb.PostCode) error {
	// TODO: impl & check format
	client := ctx.Value(meta.KeyRepoRedis).(*redis.Client)
	if msg.Type == pb.PostType_CodePhone ||
		msg.Type == pb.PostType_CodeEmail {
		code := newCode()

		k := KeyGroupRedisVerifyCodes[msg.Type] + msg.To
		client.Set(k, code, initialize.TimeVerifyCode)

		return nil
	}
	return meta.ErrUnsupportedVerifyType
}
