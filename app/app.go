package app

import (
	"context"

	api "github.com/MisLink/go-web-template/api"
	"github.com/MisLink/go-web-template/pkg/crontab"
	"github.com/MisLink/go-web-template/pkg/server"

	"github.com/getsentry/sentry-go"
	"github.com/google/wire"
	"github.com/twitchtv/twirp"
)

type Server struct{}

func NewServer() api.TwirpServer {
	s := &Server{}
	handler := api.NewMODULE_NAMEServer(
		s,
		twirp.WithServerPathPrefix(""),
		twirp.WithServerHooks(&twirp.ServerHooks{
			Error: func(ctx context.Context, err twirp.Error) context.Context {
				if hub := sentry.GetHubFromContext(ctx); hub != nil {
					hub.CaptureException(err)
				}
				return ctx
			},
		}))
	return handler
}

func NewHandler() server.MountFunc {
	return func(m server.Mount) {
		m.Mount("/", NewServer())
	}
}

func NewCrontab() crontab.RegisterFunc {
	return func(r crontab.Register) {}
}

var ProviderSet = wire.NewSet(NewHandler, NewCrontab)
