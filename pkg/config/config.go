package config

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/wire"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	DSN string `koanf:"dsn"`
	ENV string `koanf:"env"`
}

func NewConfig(k *koanf.Koanf) (*Config, error) {
	c := new(Config)
	if err := k.Unmarshal("", c); err != nil {
		return nil, err
	}
	return c, nil
}

func New() (*koanf.Koanf, error) {
	_, filename, _, _ := runtime.Caller(0) // nolint:dogsled
	dir := filepath.Join(path.Dir(filename), "..", "..")
	_ = godotenv.Load(filepath.Join(dir, ".env"))
	k := koanf.New(".")
	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "__", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	return k, nil
}

var ProviderSet = wire.NewSet(New, NewConfig)
