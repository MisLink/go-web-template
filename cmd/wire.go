//go:build wireinject

package cmd

import (
	"github.com/MisLink/go-web-template/app"

	"github.com/MisLink/go-web-template/pkg/config"
	"github.com/MisLink/go-web-template/pkg/crontab"
	"github.com/MisLink/go-web-template/pkg/database"
	"github.com/MisLink/go-web-template/pkg/database/ent"
	"github.com/MisLink/go-web-template/pkg/logger"
	"github.com/MisLink/go-web-template/pkg/redis"
	"github.com/MisLink/go-web-template/pkg/server"
	"github.com/MisLink/go-web-template/pkg/telemetry"

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
