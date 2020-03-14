package initialize

import (
	"github.com/ghosv/open/plat/conf"

	"github.com/go-redis/redis/v7"
)

// Redis init
func Redis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.ENV.Repo.Redis.URL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
