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
	s3Repo "github.com/leftovers-2025/kds_backend/internal/kds/repository/s3"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

var datasourceSet = wire.NewSet(
	datasource.NewGoogleApiClient,
	datasource.GetMySqlConnection,
	datasource.GetRedisClient,
	datasource.GetMinIOClient,
)

var repositorySet = wire.NewSet(
	apiRepo.NewApiGoogleRepository,
	mysqlRepo.NewMySqlUserRepository,
	redisRepo.NewRedisTokenRepository,
	mysqlRepo.NewMySqlTagRepository,
	mysqlRepo.NewMySqlLocationRepository,
	mysqlRepo.NewMySqlPostRepository,
	s3Repo.NewS3Repository,
)

var serviceSet = wire.NewSet(
	service.NewUserCommandService,
	service.NewGoogleCommandService,
	service.NewAuthQueryService,
	service.NewAuthCommandService,
	service.NewUserQueryService,
	service.NewTagQueryService,
	service.NewTagCommandService,
	service.NewLocationQueryService,
	service.NewLocationCommandService,
	service.NewPostCommandService,
)

var handlerSet = wire.NewSet(
	handler.NewGoogleHandler,
	handler.NewErrorHandler,
	handler.NewUserHandler,
	handler.NewAuthHandler,
	handler.NewTagHandler,
	handler.NewLocationHandler,
	handler.NewPostHandler,
)

type HandlerSets struct {
	GoogleHandler   *handler.GoogleHandler
	ErrorHandler    *handler.ErrorHandler
	AuthHandler     *handler.AuthHandler
	UserHandler     *handler.UserHandler
	TagHandler      *handler.TagHandler
	LocationHandler *handler.LocationHandler
	PostHandler     *handler.PostHandler
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
