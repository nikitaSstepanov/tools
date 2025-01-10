package tools

import (
	"github.com/nikitaSstepanov/tools/client/mail"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/httper"
	"github.com/nikitaSstepanov/tools/sl"
	"github.com/nikitaSstepanov/tools/utils/coder"
)

type pgConfig struct {
	Postgres pg.Config `yaml:"postgres"`
}

type redisConfig struct {
	Redis redis.Config `yaml:"redis"`
}

type slConfig struct {
	Logger sl.Config `yaml:"logger"`
}

type httpServerConfig struct {
	Server httper.ServerCfg `yaml:"server"`
}

type mailConfig struct {
	Mail mail.Config `yaml:"mail"`
}

type coderConfig struct {
	Coder coder.Config `yaml:"coder"`
}
