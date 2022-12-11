package utils

import (
	"context"
	"os/signal"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

func GracefulStop(
	ctx context.Context,
	startFunc func(ctx context.Context) error,
	stopFunc func() error,
) error {
	ctx, stop := signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT, unix.SIGHUP, unix.SIGQUIT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return startFunc(ctx) })
	eg.Go(func() error {
		<-ctx.Done()
		return stopFunc()
	})
	return eg.Wait()
}
