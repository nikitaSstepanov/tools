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
)

func Init(path ...string) error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	cfgPath := os.Getenv("CONFIG_PATH")

	if cfgPath != "" {
		configPath = cfgPath
		return nil
	}

	if len(path) != 0 {
		if len(path) > 1 {
			return errors.New("there should be only one config path")
		}

		configPath = path[0]
		return nil
	}

	return nil
}

func Pg(path ...string) (pg.Client, error) {
	cfgPath := configPath

	if len(path) != 0 {
		if len(path) > 1 {
			return nil, errors.New("there should be only one config path")
		}

		cfgPath = path[0]
	}

	var cfg pgConfig

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	ctx := context.Background()

	postgres, err := pg.New(ctx, &cfg.Postgres)
	if err != nil {
		return nil, err
	}

	if cfg.Postgres.MigrationsRun {
		err := migrate.MigratePg(postgres.ToPgx(), cfg.Postgres.MigrationsPath)
		if err != nil {
			return nil, err
		}
	}

	return postgres, nil
}

func Redis(path ...string) (redis.Client, error) {
	cfgPath := configPath

	if len(path) != 0 {
		if len(path) > 1 {
			return redis.Client{}, errors.New("there should be only one config path")
		}

		cfgPath = path[0]
	}

	var cfg redisConfig

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return redis.Client{}, err
	}

	ctx := context.Background()

	rs, err := redis.New(ctx, &cfg.Redis)
	if err != nil {
		return redis.Client{}, err
	}

	return rs, nil
}

func Sl(path ...string) (*sl.Logger, error) {
	cfgPath := configPath

	if len(path) != 0 {
		if len(path) > 1 {
			return nil, errors.New("there should be only one config path")
		}

		cfgPath = path[0]
	}

	var cfg slConfig

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	return sl.New(&cfg.Logger), nil
}

func HttpServer(handler http.Handler, path ...string) (*httper.Server, error) {
	cfgPath := configPath

	if len(path) != 0 {
		if len(path) > 1 {
			return nil, errors.New("there should be only one config path")
		}

		cfgPath = path[0]
	}

	var cfg httpServerConfig

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	return httper.NewServer(&cfg.Server, handler), nil
}

func Mail(path ...string) (*mail.Client, error) {
	cfgPath := configPath

	if len(path) != 0 {
		if len(path) > 1 {
			return nil, errors.New("there should be only one config path")
		}

		cfgPath = path[0]
	}

	var cfg mailConfig

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	return mail.New(&cfg.Mail), nil
}

func Coder(path ...string) (*coder.Coder, error) {
	cfgPath := configPath

	if len(path) != 0 {
		if len(path) > 1 {
			return nil, errors.New("there should be only one config path")
		}

		cfgPath = path[0]
	}

	var cfg coderConfig

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	return coder.New(&cfg.Coder), nil
}
