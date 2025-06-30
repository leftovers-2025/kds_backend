package handler

import "github.com/labstack/echo/v4"

type GoogleHandler struct {
}

func NewGoogleHandler() *GoogleHandler {
	return &GoogleHandler{}
}

func (h *GoogleHandler) Redirect(ctx echo.Context) error {
	return nil
}

func (h *GoogleHandler) Code(ctx echo.Context) error {
	return nil
}
