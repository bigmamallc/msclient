package msclient

import (
	"github.com/bigmamallc/env"
	"time"
)

type Cfg struct {
	BaseURL        string        `env:"BASE_URL" required:"true"`
	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT" default:"5s"`
	MaxRetries     int           `env:"MAX_RETRIES" default:"3"`

	BasicAuthUsername string `env:"BASIC_AUTH_USERNAME" default:""`
	BasicAuthPassword string `env:"BASIC_AUTH_PASSWORD" default:""`
}

func EnvCfg(envPrefix string) (*Cfg, error) {
	c := &Cfg{}
	if err := env.SetWithEnvPrefix(c, envPrefix); err != nil {
		return nil, err
	}
	return c, nil
}

func MustEnvCfg(envPrefix string) *Cfg {
	c, err := EnvCfg(envPrefix)
	if err != nil {
		panic(err)
	}
	return c
}
