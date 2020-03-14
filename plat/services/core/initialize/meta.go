package initialize

import (
	"time"

	"github.com/ghosv/open/meta"
)

// meta
const (
	SrvName = meta.SrvCore
	SrvVer  = meta.Version

	// TODO: 重新考虑过期时间、刷新时间
	DefaultJWTTerm  = time.Hour * 24 // 过期时间
	DefaultJWTFresh = time.Hour * 6  // 刷新时间

	// TODO: 重新考虑授权码有效期
	DefaultAuthCodeTime = time.Minute * 10
)
