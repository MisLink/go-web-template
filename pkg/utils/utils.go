package utils

import (
	"context"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func Lifecycle(ctx context.Context, start func() error, stop func() error) error {
	ctx, reset := signal.NotifyContext(
		ctx,
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	defer reset()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return start() })
	eg.Go(func() error {
		<-ctx.Done()
		return stop()
	})
	return eg.Wait()
}
