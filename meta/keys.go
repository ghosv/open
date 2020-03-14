package meta

type key struct {
	payload string
}

// key of context
var (
	// === gate ===
	KeyService      = key{CtxService}
	KeyTokenPayload = key{CtxTokenPayload}

	// === services/repo ===
	KeyRepo      = key{"key.repo.mysql"} // MySQL
	KeyRepoRedis = key{"key.repo.redis"}
	// KeyRepoMongo = struct{}{}

	// === client ===
	KeyClientNotifyVerify = key{"client.notify.verify"}
)

// string key
const (
	// === gate/sola ===
	CtxService      = "open.service"
	CtxTokenPayload = "open.token.payload"

	// === redis ===
	KeyRedisPhoneCode   = "verifycode.phone-"
	KeyRedisEmailCode   = "verifycode.email-"
	KeyRedisGhostIMCode = "verifycode.ghostim-"
)
