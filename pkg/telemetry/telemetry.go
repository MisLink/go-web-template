package telemetry

import (
	"context"

	"MODULE_NAME/types"

	"github.com/google/wire"
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

func NewResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(types.ModuleName),
			semconv.ServiceVersion(types.Version),
			semconv.DeploymentEnvironment(""),
		),
	)
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
	return tp, func() { _ = tp.Shutdown(context.Background()) }, nil
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
	return mp, func() { _ = mp.Shutdown(context.Background()) }, nil
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
