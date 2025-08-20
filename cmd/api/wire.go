//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/leftovers-2025/kds_backend/internal/kds/datasource"
	"github.com/leftovers-2025/kds_backend/internal/kds/handler"
	apiRepo "github.com/leftovers-2025/kds_backend/internal/kds/repository/api"
	emailRepo "github.com/leftovers-2025/kds_backend/internal/kds/repository/email"
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
	emailRepo.EmailFromEnv,
)

var repositorySet = wire.NewSet(
	apiRepo.NewApiGoogleRepository,
	mysqlRepo.NewMySqlUserRepository,
	redisRepo.NewRedisTokenRepository,
	mysqlRepo.NewMySqlTagRepository,
	mysqlRepo.NewMySqlLocationRepository,
	mysqlRepo.NewMySqlPostRepository,
	mysqlRepo.NewNotificationRepository,
	s3Repo.NewS3Repository,
	emailRepo.NewEmailRepository,
)

var serviceSet = wire.NewSet(
	service.NewGoogleCommandService,
	service.NewAuthQueryService,
	service.NewAuthCommandService,
	service.NewUserQueryService,
	service.NewUserCommandService,
	service.NewUserEditCommandService,
	service.NewTagQueryService,
	service.NewTagCommandService,
	service.NewLocationQueryService,
	service.NewLocationCommandService,
	service.NewPostQueryService,
	service.NewPostCommandService,
	service.NewNotificationCommandService,
)

var handlerSet = wire.NewSet(
	handler.NewGoogleHandler,
	handler.NewErrorHandler,
	handler.NewUserHandler,
	handler.NewAuthHandler,
	handler.NewTagHandler,
	handler.NewLocationHandler,
	handler.NewPostHandler,
	handler.NewNotificationHandler,
)

type HandlerSets struct {
	GoogleHandler       *handler.GoogleHandler
	ErrorHandler        *handler.ErrorHandler
	AuthHandler         *handler.AuthHandler
	UserHandler         *handler.UserHandler
	TagHandler          *handler.TagHandler
	LocationHandler     *handler.LocationHandler
	PostHandler         *handler.PostHandler
	NotificationHandler *handler.NotificationHandler
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
