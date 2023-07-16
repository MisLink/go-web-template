package database

import (
	"context"
	"fmt"

	"MODULE_NAME/pkg/database/ent"

	_ "MODULE_NAME/pkg/database/ent/runtime"

	"entgo.io/ent/dialect/sql"
	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type Options struct {
	Driver string
	URL    string
	Debug  bool
}

func NewOptions(k *koanf.Koanf) (*Options, error) {
	o := new(Options)
	if err := k.Unmarshal("database", o); err != nil {
		return nil, err
	}
	return o, nil
}

func New(opt *Options, logger zerolog.Logger, tp trace.TracerProvider, mp metric.MeterProvider) (*ent.Client, error) {
	db, err := otelsql.Open(opt.Driver, opt.URL,
		otelsql.WithAttributes(semconv.DBSystemMySQL), otelsql.WithTracerProvider(tp), otelsql.WithMeterProvider(mp))
	if err != nil {
		return nil, err
	}
	if err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemMySQL)); err != nil {
		return nil, err
	}
	drv := sql.OpenDB(opt.Driver, db)
	logger = logger.With().Str("logger", "database").Logger()
	opts := []ent.Option{ent.Driver(drv), ent.Log(func(a ...any) {
		logger.Print(a...)
	})}
	if opt.Debug {
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
			err = fmt.Errorf("%w: rolling back transaction: %w", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
