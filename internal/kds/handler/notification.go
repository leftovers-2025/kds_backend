package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

var (
	ErrNotificationInvalidTagId      = common.NewValidationError(errors.New("invalid tag ids"))
	ErrNotificationInvalidLocationId = common.NewValidationError(errors.New("invalid location ids"))
)

type NotificationHandler struct {
	notificationCmdService *service.NotificationCommandService
}

func NewNotificationHandler(
	notificationCmdService *service.NotificationCommandService,
) *NotificationHandler {
	if notificationCmdService == nil {
		panic("nil NotificationCommandService")
	}
	return &NotificationHandler{
		notificationCmdService: notificationCmdService,
	}
}

type NotificationSettingsRequest struct {
	Enabled     bool     `json:"enabled"`
	TagIds      []string `json:"tagIds"`
	LocationIds []string `json:"locationIds"`
}

// SaveSettings godoc
//
//	@Summary		Save notification settings
//	@Description	Save notification settings for the user
//	@Tags			notifications
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	NotificationSettingsRequest	true	"Notification settings request"
//	@Success		204
//	@Router			/notifications [put]
//
// 通知設定を保存する
func (h *NotificationHandler) SaveSettings(ctx echo.Context) error {
	// ユーザーId取得
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return err
	}
	// 通知保存リクエスト取得
	request := NotificationSettingsRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	// タグIdをUUID変換
	tagIds := []uuid.UUID{}
	for _, tagId := range request.TagIds {
		id, err := uuid.Parse(tagId)
		if err != nil {
			return ErrNotificationInvalidTagId
		}
		tagIds = append(tagIds, id)
	}
	// ロケーションIdをUUID変換
	locationIds := []uuid.UUID{}
	for _, locationId := range request.LocationIds {
		id, err := uuid.Parse(locationId)
		if err != nil {
			return ErrNotificationInvalidLocationId
		}
		locationIds = append(locationIds, id)
	}
	err = h.notificationCmdService.SaveNotification(service.NotificationCommandInput{
		UserId:      userId,
		Enabled:     request.Enabled,
		LocationIds: locationIds,
		TagIds:      tagIds,
	})
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
