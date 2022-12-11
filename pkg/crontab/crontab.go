package crontab

import (
	"context"
	"time"

	"MODULE_NAME/pkg/utils"

	"github.com/getsentry/sentry-go"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rueian/rueidis/rueidislock"
)

type PeriodicTask func() (string, func() error)

type Crontab struct {
	locker rueidislock.Locker
	cron   *cron.Cron
	logger zerolog.Logger
}

type Logger struct {
	logger  zerolog.Logger
	logInfo bool
}

func (l *Logger) Info(msg string, keysAndValues ...any) {
	if l.logInfo {
		l.logger.Info().Msgf(msg, keysAndValues...)
	}
}

func (l *Logger) Error(err error, msg string, keysAndValues ...any) {
	l.logger.Err(err).Msgf(msg, keysAndValues...)
}

type Register interface {
	Register(string, PeriodicTask) error
}

type RegisterFunc func(Register)

func New(
	logger zerolog.Logger,
	locker rueidislock.Locker,
	register RegisterFunc,
) *Crontab {
	logger = logger.With().Str("logger", "crontab").Logger()
	cronLogger := &Logger{
		logger:  logger,
		logInfo: true,
	}
	tz, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		tz = time.Local
	}
	crontab := &Crontab{
		cron: cron.New(
			cron.WithLogger(cronLogger),
			cron.WithLocation(tz),
			cron.WithChain(cron.Recover(cronLogger)),
		),
		logger: logger,
		locker: locker,
	}
	register(crontab)
	return crontab
}

func (c *Crontab) Register(name string, task PeriodicTask) error {
	spec, fn := task()
	logger := c.logger.With().Str("task", name).Logger()
	hub := sentry.CurrentHub()
	hub.Scope().SetTransaction(name)
	_, err := c.cron.AddFunc(spec, func() {
		logger.Debug().Msg("crontab executing")
		err := fn()
		if err != nil {
			logger.Err(err).Stack().Msg("execute error")
			hub.CaptureException(err)
		} else {
			logger.Debug().Msg("execute success")
		}
	})
	return err
}

func (c *Crontab) Start() error {
	ctx := context.Background()
	for {
		ctx, cancel, err := c.locker.WithContext(ctx, "crontab")
		if err != nil {
			c.logger.Err(err).Msg("obtain lock error")
			continue
		}
		defer cancel()
		return utils.GracefulStop(ctx, func(ctx context.Context) error {
			c.cron.Run()
			return nil
		}, func() error {
			ctx := c.cron.Stop()
			<-ctx.Done()
			return nil
		})
	}
}

var ProviderSet = wire.NewSet(New)
