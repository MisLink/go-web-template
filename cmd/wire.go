//go:build wireinject

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

var providerSet = wire.NewSet(
	config.ProviderSet,
	telemetry.ProviderSet,
	logger.ProviderSet,
	server.ProviderSet,
	app.ProviderSet,
	crontab.ProviderSet,
	redis.ProviderSet,
	database.ProviderSet,
)

func CreateServer() (*server.Server, func(), error) {
	panic(wire.Build(providerSet))
}

func CreateCrontab() (*crontab.Crontab, func(), error) {
	panic(wire.Build(providerSet))
}

func CreateDatabase() (*ent.Client, func(), error) {
	panic(wire.Build(providerSet))
}
