package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(".env is not found")
	}

	e := echo.New()
	e.Use(middleware.Logger())

	handlerSets := InitHandlerSets()

	// error handling
	e.HTTPErrorHandler = handlerSets.ErrorHandler.HandleError

	// google oauth
	e.GET("/oauth/google/redirect", handlerSets.GoogleHandler.Redirect)

	e.Logger.Fatal(e.Start(":8630"))
}
