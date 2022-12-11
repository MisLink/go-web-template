package database

import (
	"MODULE_NAME/pkg/database/ent"

	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/knadh/koanf"
	"github.com/rs/zerolog"
)

type DBConfig struct {
	Driver string
	Url    string
	Debug  bool
}

type Options struct {
	Default DBConfig
}

func NewOptions(k *koanf.Koanf) (*Options, error) {
	o := new(Options)
	if err := k.Unmarshal("database", o); err != nil {
		return nil, err
	}
	return o, nil
}

func New(opt *Options, logger zerolog.Logger) (*ent.Client, error) {
	drv, err := sql.Open(opt.Default.Driver, opt.Default.Url)
	if err != nil {
		return nil, err
	}
	logger = logger.With().Str("logger", "database").Logger()
	opts := []ent.Option{ent.Driver(drv), ent.Log(func(a ...any) {
		logger.Print(a...)
	})}
	if opt.Default.Debug {
		opts = append(opts, ent.Debug())
	}
	client := ent.NewClient(opts...)
	return client, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
