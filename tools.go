package tools

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/nikitaSstepanov/tools/client/mail"
	"github.com/nikitaSstepanov/tools/client/pg"
	"github.com/nikitaSstepanov/tools/client/redis"
	"github.com/nikitaSstepanov/tools/httper"
	"github.com/nikitaSstepanov/tools/migrate"
	"github.com/nikitaSstepanov/tools/sl"
	"github.com/nikitaSstepanov/tools/utils/coder"
)

var (
	configPath = "config/config.yaml"
	config     = &toolsConfig{}
)

func Init(path ...string) error {
	if len(path) > 1 {
		return errors.New("there should be only one config path")
	}

	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	cfgPath := os.Getenv("CONFIG_PATH")

	if cfgPath != "" {
		configPath = cfgPath
	} else if len(path) != 0 {
		configPath = path[0]
	}

	var cfg toolsConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil
	}

	config = &cfg

	return nil
}

func Pg() (pg.Client, error) {
	ctx := context.Background()

	postgres, err := pg.New(ctx, &config.Postgres)
	if err != nil {
		return nil, err
	}

	if config.Postgres.MigrationsRun {
		err := migrate.MigratePg(postgres.ToPgx(), config.Postgres.MigrationsPath)
		if err != nil {
			return nil, err
		}
	}

	return postgres, nil
}

func Redis() (redis.Client, error) {
	ctx := context.Background()

	rs, err := redis.New(ctx, &config.Redis)
	if err != nil {
		return redis.Client{}, err
	}

	return rs, nil
}

func Sl() *sl.Logger {
	return sl.New(&config.Logger)
}

func HttpServer(handler http.Handler) *httper.Server {
	return httper.NewServer(&config.HttpServer, handler)
}

func Mail() *mail.Client {
	return mail.New(&config.Mail)
}

func Coder() *coder.Coder {
	return coder.New(&config.Coder)
}
