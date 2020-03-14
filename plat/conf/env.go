package conf

import (
	"github.com/micro/go-micro/v2/util/log"

	"github.com/timest/env"
)

type config struct {
	Repo struct {
		Mongo struct {
			URL string `env:"URL" default:"root:123456@127.0.0.1:27017"`
		} `env:"MONGO"`
		Redis struct {
			URL string `env:"URL" default:"127.0.0.1:6379"`
		} `env:"REDIS"`
		MySQL struct {
			Host string `env:"HOST" default:"127.0.0.1:3306"`
			User string `env:"USER" default:"root"`
			Pass string `env:"PASS" default:"123456"`
		} `env:"MYSQL"`
	} `env:"REPO"`
	Key struct {
		JWT string `env:"JWT" default:"jwt_key"`
	} `env:"KEY"`
}

func (c *config) KeyJWT() []byte {
	return []byte(c.Key.JWT)
}

// config(s)
var (
	ENV = new(config)
)

func init() {
	if err := env.Fill(ENV); err != nil {
		log.Fatal(err)
	}
}
