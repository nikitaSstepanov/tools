package tools

import (
	"github.com/nikitaSstepanov/tools/client/mail"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/httper"
	"github.com/nikitaSstepanov/tools/sl"
	"github.com/nikitaSstepanov/tools/utils/coder"
)

type toolsConfig struct {
	Postgres   pg.Config        `yaml:"postgres"`
	Redis      redis.Config     `yaml:"redis"`
	Logger     sl.Config        `yaml:"logger"`
	HttpServer httper.ServerCfg `yaml:"http_server"`
	Mail       mail.Config      `yaml:"mail"`
	Coder      coder.Config     `yaml:"coder"`
}
