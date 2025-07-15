package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

type UserHandler struct {
	userQueryService   *service.UserQueryService
	userEditCmdService *service.UserEditCommandService
}

func NewUserHandler(
	userQueryService *service.UserQueryService,
	userEditCmdService *service.UserEditCommandService,
) *UserHandler {
	if userQueryService == nil {
		panic("nil UserQueryService")
	}
	if userEditCmdService == nil {
		panic("nil UserEditCommandService")
	}
	return &UserHandler{
		userQueryService:   userQueryService,
		userEditCmdService: userEditCmdService,
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

type UserEditRequest struct {
	Role string `json:"role"`
}

// ユーザーを編集
func (h *UserHandler) EditUser(ctx echo.Context) error {
	// ユーザーId取得
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return err
	}
	// 編集対象ユーザーId取得
	targetUserId, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		return err
	}
	// ユーザー編集リクエスト取得
	request := UserEditRequest{}
	if err = ctx.Bind(&request); err != nil {
		return common.NewValidationError(err)
	}
	// ユーザーを編集
	err = h.userEditCmdService.EditRole(userId, service.UserEditRoleCommandInput{
		TargetUserId: targetUserId,
		Role:         request.Role,
	})
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
