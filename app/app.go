package app

import (
	"context"

	"MODULE_NAME/api"
	"MODULE_NAME/pkg/crontab"
	"MODULE_NAME/pkg/server"

	"github.com/getsentry/sentry-go"
	"github.com/google/wire"
	"github.com/twitchtv/twirp"
)

type Server struct{}

func NewServer() api.TwirpServer {
	server := &Server{}
	handler := api.NewMODULE_NAMEServer(
		server,
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
	return func(m server.Mounter) {
		m.Mount("/", NewServer())
	}
}

func NewCrontab() crontab.RegisterFunc {
	return func(r crontab.Register) {}
}

var ProviderSet = wire.NewSet(NewHandler, NewCrontab)
