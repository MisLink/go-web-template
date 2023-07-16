// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cmd

import (
	"MODULE_NAME/app"
	"MODULE_NAME/pkg/config"
	"MODULE_NAME/pkg/crontab"
	"MODULE_NAME/pkg/database"
	"MODULE_NAME/pkg/database/ent"
	"MODULE_NAME/pkg/logger"
	"MODULE_NAME/pkg/redis"
	"MODULE_NAME/pkg/server"
	"MODULE_NAME/pkg/telemetry"
	"github.com/google/wire"
)

// Injectors from wire.go:

func CreateServer() (*server.Server, func(), error) {
	koanf, cleanup, err := config.New()
	if err != nil {
		return nil, nil, err
	}
	options, err := server.NewOptions(koanf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	mountFunc := app.NewHandler()
	loggerOptions, err := logger.NewOptions(koanf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	zerologLogger := logger.New(loggerOptions)
	resource, err := telemetry.NewResource()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	spanExporter, err := telemetry.NewTraceExporter()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	tracerProvider, cleanup2, err := telemetry.NewTraceProvider(resource, spanExporter)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	reader, err := telemetry.NewMetricReader()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	meterProvider, cleanup3, err := telemetry.NewMeterProvider(resource, reader)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	serverServer, err := server.New(options, mountFunc, zerologLogger, tracerProvider, meterProvider)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return serverServer, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

func CreateCrontab() (*crontab.Crontab, func(), error) {
	koanf, cleanup, err := config.New()
	if err != nil {
		return nil, nil, err
	}
	options, err := logger.NewOptions(koanf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	zerologLogger := logger.New(options)
	redisOptions, err := redis.NewOptions(koanf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	clientOption, err := redis.NewRedisOption(redisOptions)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	locker, err := redis.NewLock(clientOption, redisOptions)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	registerFunc := app.NewCrontab()
	crontabCrontab := crontab.New(zerologLogger, locker, registerFunc)
	return crontabCrontab, func() {
		cleanup()
	}, nil
}

func CreateDatabase() (*ent.Client, func(), error) {
	koanf, cleanup, err := config.New()
	if err != nil {
		return nil, nil, err
	}
	options, err := database.NewOptions(koanf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	loggerOptions, err := logger.NewOptions(koanf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	zerologLogger := logger.New(loggerOptions)
	resource, err := telemetry.NewResource()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	spanExporter, err := telemetry.NewTraceExporter()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	tracerProvider, cleanup2, err := telemetry.NewTraceProvider(resource, spanExporter)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	reader, err := telemetry.NewMetricReader()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	meterProvider, cleanup3, err := telemetry.NewMeterProvider(resource, reader)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	client, err := database.New(options, zerologLogger, tracerProvider, meterProvider)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return client, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var providerSet = wire.NewSet(config.ProviderSet, telemetry.ProviderSet, logger.ProviderSet, server.ProviderSet, app.ProviderSet, crontab.ProviderSet, redis.ProviderSet, database.ProviderSet)
