package config

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"MODULE_NAME/types"

	"github.com/getsentry/sentry-go"
	"github.com/google/wire"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

func New() (*koanf.Koanf, func(), error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(path.Dir(filename), "..", "..")
	_ = godotenv.Load(filepath.Join(dir, ".env"))
	k := koanf.New(".")
	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, nil, err
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:            k.String("dsn"),
		Environment:    k.String("env"),
		SampleRate:     0.1,
		SendDefaultPII: true,
		Release:        types.Version,
	}); err != nil {
		return nil, nil, err
	}
	return k, func() { sentry.Flush(time.Second * 5) }, nil
}

var ProviderSet = wire.NewSet(New)
