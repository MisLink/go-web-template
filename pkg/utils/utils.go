package utils

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func Lifecycle(ctx context.Context, start func(context.Context) error, stop func() error) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return start(ctx) })
	eg.Go(func() error {
		<-ctx.Done()
		return stop()
	})
	return eg.Wait()
}
