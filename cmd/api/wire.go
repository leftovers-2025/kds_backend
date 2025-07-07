//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/leftovers-2025/kds_backend/internal/kds/datasource"
	"github.com/leftovers-2025/kds_backend/internal/kds/handler"
	apiRepo "github.com/leftovers-2025/kds_backend/internal/kds/repository/api"
	mysqlRepo "github.com/leftovers-2025/kds_backend/internal/kds/repository/mysql"
	redisRepo "github.com/leftovers-2025/kds_backend/internal/kds/repository/redis"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

var datasourceSet = wire.NewSet(
	datasource.NewGoogleApiClient,
	datasource.GetMySqlConnection,
	datasource.GetRedisClient,
)

var repositorySet = wire.NewSet(
	apiRepo.NewApiGoogleRepository,
	mysqlRepo.NewMySqlUserRepository,
	redisRepo.NewRedisTokenRepository,
)

var serviceSet = wire.NewSet(
	service.NewUserCommandService,
	service.NewGoogleCommandService,
	service.NewAuthQueryService,
	service.NewAuthCommandService,
	service.NewUserQueryService,
)

var handlerSet = wire.NewSet(
	handler.NewGoogleHandler,
	handler.NewErrorHandler,
	handler.NewUserHandler,
	handler.NewAuthHandler,
)

type HandlerSets struct {
	GoogleHandler *handler.GoogleHandler
	ErrorHandler  *handler.ErrorHandler
	AuthHandler   *handler.AuthHandler
	UserHandler   *handler.UserHandler
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
