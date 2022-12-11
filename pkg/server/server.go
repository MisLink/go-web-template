package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"MODULE_NAME/pkg/utils"
	"MODULE_NAME/types"

	_ "MODULE_NAME/pkg/metrics"

	"code.cloudfoundry.org/bytefmt"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/wire"
	"github.com/knadh/koanf"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Options struct {
	Addr string
	Dsn  string
	Env  string
}

func NewOptions(k *koanf.Koanf) (*Options, error) {
	o := new(Options)
	if err := k.Unmarshal("server", o); err != nil {
		return nil, err
	}
	return o, nil
}

type Server struct {
	server http.Server
	logger zerolog.Logger
}

type Mounter interface {
	Mount(string, http.Handler)
}

type MountFunc func(Mounter)

func New(
	opt *Options,
	mount MountFunc,
	logger zerolog.Logger,
) (*Server, error) {
	logger = logger.With().Str("logger", "server").Logger()
	router := chi.NewRouter()
	router.Use(
		sentryhttp.New(sentryhttp.Options{Repanic: true}).Handle,
		hlog.NewHandler(logger),
		middleware.RealIP,
		hlog.RequestIDHandler("request_id", "X-Request-Id"),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Int("status", status).
				Int("size", size).
				Str("method", r.Method).
				Stringer("url", r.URL).
				Dur("duration", duration).
				Str("ip", r.RemoteAddr).
				Str("user_agent", r.Header.Get("User-Agent")).
				Str("referer", r.Header.Get("Referer")).
				Str("proto", r.Proto).
				Msg("")
		}),
		middleware.CleanPath,
		middleware.Heartbeat("/ping"),
		middleware.Recoverer,
	)
	router.Group(func(r chi.Router) {
		r.Use(
			middleware.BasicAuth("MODULE_NAME-debug", map[string]string{"MODULE_NAME": ""}),
		)
		r.Mount("/debug", middleware.Profiler())
	})
	router.Mount("/metrics", promhttp.Handler())
	mount(router)
	server := &Server{
		server: http.Server{
			Addr:              opt.Addr,
			Handler:           router,
			ReadTimeout:       60 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      60 * time.Second,
			MaxHeaderBytes:    1 * bytefmt.MEGABYTE,
		},
		logger: logger,
	}
	return server, nil
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}
	s.logger.Info().
		Str("addr", s.server.Addr).
		Str("version", types.Version).
		Str("built_at", types.BuiltAt).
		Msg("start listening")
	return utils.GracefulStop(
		context.Background(),
		func(ctx context.Context) error { return s.server.Serve(ln) },
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			return s.server.Shutdown(ctx)
		},
	)
}

var ProviderSet = wire.NewSet(NewOptions, New)
