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

	handlerSets := InitHandlerSets()
	auth := e.Group("", handlerSets.AuthHandler.JwtAuthorization)

	// error handling
	e.HTTPErrorHandler = handlerSets.ErrorHandler.HandleError

	// google oauth
	e.GET("/oauth/google/redirect", handlerSets.GoogleHandler.Redirect)

	// user
	auth.GET("/users/@me", handlerSets.UserHandler.Me)

	e.Logger.Fatal(e.Start(":" + port))
}
