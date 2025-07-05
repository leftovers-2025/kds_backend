package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) HandleError(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}
	code := http.StatusInternalServerError
	message := "internal server error"
	// バリデーションエラー
	if internalErr, ok := err.(*common.ValidationError); ok {
		code = http.StatusBadRequest
		message = internalErr.Error()
	}
	ctx.Logger().Error(err)
	err = ctx.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
	})
	if err != nil {
		ctx.Logger().Error(err)
	}
}
