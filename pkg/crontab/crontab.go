package crontab

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/redis/rueidis/rueidislock"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type PeriodicTask func() (string, func() error)

type Crontab struct {
	locker rueidislock.Locker
	cron   *cron.Cron
	logger zerolog.Logger
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
	_, err := c.cron.AddFunc(spec, func() {
		logger.Debug().Msg("crontab executing")
		err := fn()
		if err != nil {
			logger.Err(err).Stack().Msg("execute error")
		} else {
			logger.Debug().Msg("execute success")
		}
	})
	return err
}

func (c *Crontab) Lock(ctx context.Context) (context.Context, context.CancelFunc, error) {
	var nCtx context.Context
	var cancel context.CancelFunc
	for {
		var err error
		nCtx, cancel, err = c.locker.WithContext(ctx, "crontab")
		if err != nil {
			c.logger.Err(err).Msg("obtain lock error")
			continue
		}
		break
	}
	return nCtx, cancel, nil
}

func (c *Crontab) Start() error {
	c.cron.Run()
	return nil
}

func (c *Crontab) Close() error {
	<-c.cron.Stop().Done()
	return nil
}

var ProviderSet = wire.NewSet(New)
