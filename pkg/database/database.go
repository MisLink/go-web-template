package database

import (
	"context"
	"fmt"

	"MODULE_NAME/pkg/database/ent"

	_ "MODULE_NAME/pkg/database/ent/runtime"

	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type DBConfig struct {
	Driver string
	Url    string // revive:disable-line
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

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
