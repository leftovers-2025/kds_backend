//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/leftovers-2025/kds_backend/internal/kds/datasource"
	"github.com/leftovers-2025/kds_backend/internal/kds/handler"
	apiRepo "github.com/leftovers-2025/kds_backend/internal/kds/repository/api"
	mysqlRepo "github.com/leftovers-2025/kds_backend/internal/kds/repository/mysql"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

var datasourceSet = wire.NewSet(
	datasource.NewGoogleApiClient,
	datasource.GetMySqlConnection,
)

var repositorySet = wire.NewSet(
	apiRepo.NewApiGoogleRepository,
	mysqlRepo.NewMySqlUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewUserCommandService,
	service.NewGoogleCommandService,
)

var handlerSet = wire.NewSet(
	handler.NewGoogleHandler,
)

type HandlerSets struct {
	GoogleHandler *handler.GoogleHandler
}

func InitHandlerSets() *HandlerSets {
	wire.Build(
		datasourceSet,
		repositorySet,
		serviceSet,
		handlerSet,
		wire.Struct(new(HandlerSets), "*"),
	)
	return nil
}
