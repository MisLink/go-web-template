package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/MisLink/go-web-template/pkg/utils"
	types "github.com/MisLink/go-web-template/types"

	"code.cloudfoundry.org/bytefmt"
	sentryhttp "github.com/getsentry/sentry-go/http"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
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
	logger zerolog.Logger
	server http.Server
}

type Mount interface {
	Mount(string, http.Handler)
}

type MountFunc func(Mount)

func New(
	opt *Options,
	mount MountFunc,
	logger zerolog.Logger,
	tp trace.TracerProvider,
	mp metric.MeterProvider,
) (*Server, error) {
	router := chi.NewRouter()
	router.Use(
		middleware.RealIP,
		middleware.CleanPath,
		hlog.NewHandler(logger),
		sentryhttp.New(sentryhttp.Options{Repanic: true}).Handle,
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
				Send()
		}),
		middleware.Heartbeat("/ping"),
		middleware.Recoverer,
	)
	router.Group(func(r chi.Router) {
		r.Use(
			middleware.BasicAuth(types.ModuleName, map[string]string{types.ModuleName: ""}),
		)
		r.Mount("/debug", middleware.Profiler())
	})
	router.Mount("/metrics", promhttp.Handler())
	mount(router)
	handler := otelhttp.NewHandler(router, types.ModuleName,
		otelhttp.WithTracerProvider(tp),
		otelhttp.WithMeterProvider(mp),
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)
	server := &Server{
		server: http.Server{
			Addr:              opt.Addr,
			Handler:           handler,
			ReadTimeout:       60 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      60 * time.Second,
			MaxHeaderBytes:    1 * bytefmt.MEGABYTE,
		},
		logger: logger.With().Str("logger", "server").Logger(),
	}
	return server, nil
}

func (s *Server) Start() error {
	return utils.Lifecycle(
		context.Background(),
		func() error {
			ln, err := net.Listen("tcp", s.server.Addr)
			if err != nil {
				return err
			}
			s.logger.Info().Str("addr", s.server.Addr).Msg("start listening")
			return s.server.Serve(ln)
		},
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			return s.server.Shutdown(ctx)
		})
}

var ProviderSet = wire.NewSet(NewOptions, New)
