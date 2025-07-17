package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

var (
	ErrUserGetInvalidPage  = common.NewValidationError(errors.New("invalid page"))
	ErrUserGetInvalidLimit = common.NewValidationError(errors.New("invalid limit"))
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

// ユーザーを一覧取得する
func (h *UserHandler) GetAll(ctx echo.Context) error {
	// ユーザーID取得
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return err
	}
	// パラメーター取得
	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		return ErrPostGetInvalidLimit
	}
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return ErrPostGetInvalidPage
	}
	// ユーザーを一覧取得
	outputList, err := h.userQueryService.GetUsers(userId, service.UserAllQueryInput{
		Limit: uint(limit),
		Page:  uint(page),
	})
	if err != nil {
		return err
	}
	// レスポンスに変換
	responseList := []UserResponse{}
	for _, output := range outputList {
		responseList = append(responseList, UserResponse{
			Id:        output.Id,
			Name:      output.Name,
			Email:     output.Email,
			Role:      output.Role,
			CreatedAt: output.CreatedAt,
			UpdatedAt: output.UpdatedAt,
		})
	}
	return ctx.JSON(http.StatusOK, &responseList)
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

// ルートを譲渡
func (h *UserHandler) TransferRoot(ctx echo.Context) error {
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
	// ルート譲渡
	err = h.userEditCmdService.TransferRoot(userId, service.UserTranferRootCommandInput{
		TargetUserId: targetUserId,
	})
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
