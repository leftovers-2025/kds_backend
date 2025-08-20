package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	KDS_HTTP_PORT_DEFAULT = "8630"
)

func main() {
	// env読み込み
	err := godotenv.Load()
	if err != nil {
		panic(".env was not found")
	}

	// ポート設定
	port, ok := os.LookupEnv("KDS_HTTP_PORT")
	if !ok {
		port = KDS_HTTP_PORT_DEFAULT
	}

	e := echo.New()
	e.Use(middleware.Logger())
	api := e.Group("/api")

	handlerSets := InitHandlerSets()
	auth := api.Group("", handlerSets.AuthHandler.JwtAuthorization)

	// error handling
	e.HTTPErrorHandler = handlerSets.ErrorHandler.HandleError

	// google oauth
	e.GET("/oauth/google/redirect", handlerSets.GoogleHandler.Redirect)

	// user
	auth.GET("/users/@me", handlerSets.UserHandler.Me)
	auth.GET("/users", handlerSets.UserHandler.GetAll)
	auth.PATCH("/users/:userId/roles", handlerSets.UserHandler.EditUser)
	auth.POST("/users/:userId/root", handlerSets.UserHandler.TransferRoot)

	// tag
	api.GET("/tags", handlerSets.TagHandler.GetAll)
	auth.POST("/tags", handlerSets.TagHandler.Create)

	// location
	api.GET("/locations", handlerSets.LocationHandler.GetAll)
	auth.POST("/locations", handlerSets.LocationHandler.Create)

	// post
	api.GET("/posts", handlerSets.PostHandler.Get)
	auth.POST("/posts", handlerSets.PostHandler.Create)

	// notifications
	auth.PUT("/notifications", handlerSets.NotificationHandler.SaveSettings)

	e.Logger.Fatal(e.Start(":" + port))
}
