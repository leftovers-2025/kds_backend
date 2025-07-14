package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

type UserHandler struct {
	userQueryService *service.UserQueryService
}

func NewUserHandler(userQueryService *service.UserQueryService) *UserHandler {
	if userQueryService == nil {
		panic("nil UserQueryService")
	}
	return &UserHandler{
		userQueryService: userQueryService,
	}
}

type UserResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 自身のユーザー情報を取得する
func (h *UserHandler) Me(ctx echo.Context) error {
	// ユーザーID取得
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return err
	}
	// ユーザー情報取得
	userOutput, err := h.userQueryService.GetUser(userId)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, UserResponse{
		Id:        userOutput.Id,
		Name:      userOutput.Name,
		Email:     userOutput.Email,
		Role:      userOutput.Role,
		CreatedAt: userOutput.CreatedAt,
		UpdatedAt: userOutput.UpdatedAt,
	})
}
