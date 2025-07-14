package handler

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

type PostHandler struct {
	postCmdService *service.PostCommandService
}

func NewPostHandler(postCmdService *service.PostCommandService) *PostHandler {
	if postCmdService == nil {
		panic("nil PostCommandService")
	}
	return &PostHandler{
		postCmdService: postCmdService,
	}
}

// 投稿を新規作成
func (h *PostHandler) Create(ctx echo.Context) error {
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return err
	}
	// locationId取得
	locationId, err := uuid.Parse(ctx.FormValue("locationId"))
	if err != nil {
		return common.NewValidationError(err)
	}
	// tagIds取得
	tagIds := []uuid.UUID{}
	for tag := range strings.SplitSeq(ctx.FormValue("tagIds"), ",") {
		if strings.TrimSpace(tag) == "" {
			continue
		}
		tagId, err := uuid.Parse(strings.TrimSpace(tag))
		if err != nil {
			return common.NewValidationError(err)
		}
		tagIds = append(tagIds, tagId)
	}
	// images取得
	images := []multipart.FileHeader{}
	file, err := ctx.FormFile("image1")
	if err == nil {
		images = append(images, *file)
	}
	file, err = ctx.FormFile("image2")
	if err == nil {
		images = append(images, *file)
	}
	file, err = ctx.FormFile("image3")
	if err == nil {
		images = append(images, *file)
	}
	// 投稿作成
	err = h.postCmdService.CreatePost(userId, service.PostCreateCommandInput{
		Description: ctx.FormValue("description"),
		LocationId:  locationId,
		TagIds:      tagIds,
		Images:      images,
	})
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
