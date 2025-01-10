package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nikitaSstepanov/tools"
)

func Get(cfg interface{}) error {
	return cleanenv.ReadConfig(tools.ConfigPath, cfg)
}
