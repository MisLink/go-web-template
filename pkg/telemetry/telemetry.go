package telemetry

import (
	"context"
	"time"

	"MODULE_NAME/pkg/config"
	"MODULE_NAME/types"

	"github.com/getsentry/sentry-go"
	"github.com/google/wire"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func NewResource(c *config.Config) (*resource.Resource, func(), error) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              c.DSN,
		Environment:      c.ENV,
		SampleRate:       0.1,
		SendDefaultPII:   true,
		Release:          types.Version,
		EnableTracing:    true,
		TracesSampleRate: 0.1,
	}); err != nil {
		return nil, nil, err
	}
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(types.ModuleName),
			semconv.ServiceVersion(types.Version),
			semconv.DeploymentEnvironment(c.ENV),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	return r, func() { sentry.Flush(time.Second * 5) }, nil
}

func NewTraceExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func NewTraceProvider(r *resource.Resource, exp sdktrace.SpanExporter) (trace.TracerProvider, func(), error) {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	return tp, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown trace provider")
		}
	}, nil
}

func NewTracer(tp trace.TracerProvider) trace.Tracer {
	return tp.Tracer(types.ModuleName)
}

func NewMetricReader() (sdkmetric.Reader, error) {
	return prometheus.New(prometheus.WithNamespace(types.ModuleName))
}

func NewMeterProvider(r *resource.Resource, reader sdkmetric.Reader) (metric.MeterProvider, func(), error) {
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(r),
		sdkmetric.WithReader(reader),
	)
	otel.SetMeterProvider(mp)
	return mp, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		if err := mp.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown meter provider")
		}
	}, nil
}

func NewMeter(mp metric.MeterProvider) metric.Meter {
	return mp.Meter(types.ModuleName)
}

var ProviderSet = wire.NewSet(
	NewResource,
	NewTraceExporter,
	NewTraceProvider,
	NewTracer,
	NewMetricReader,
	NewMeterProvider,
	NewMeter,
)
