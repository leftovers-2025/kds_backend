package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

type TagHandler struct {
	tagCmdService   *service.TagCommandService
	tagQueryService *service.TagQueryService
}

func NewTagHandler(
	tagCmdService *service.TagCommandService,
	tagQueryService *service.TagQueryService,
) *TagHandler {
	if tagCmdService == nil {
		panic("nil TagCommandService")
	}
	if tagQueryService == nil {
		panic("nil TagQueryService")
	}
	return &TagHandler{
		tagCmdService:   tagCmdService,
		tagQueryService: tagQueryService,
	}
}

// タグレスポンス
type TagResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// GetAll godoc
//
//	@Summary		Get all tags
//	@Description	Get a list of all tags
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	TagResponse
//	@Router			/tags [get]
//
// タグを一覧取得
func (h *TagHandler) GetAll(ctx echo.Context) error {
	tags, err := h.tagQueryService.FindAllTags()
	if err != nil {
		return err
	}
	responseList := []TagResponse{}
	for _, tag := range tags {
		responseList = append(responseList, TagResponse{
			Id:   tag.Id.String(),
			Name: tag.Name,
		})
	}
	return ctx.JSON(http.StatusOK, &responseList)
}

// タグ作成リクエスト
type TagCreateRequest struct {
	Name string `json:"name"`
}

// Create godoc
//
//	@Summary		Create a new tag
//	@Description	Create a new tag
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	TagCreateRequest	true	"Tag create request"
//	@Success		204
//	@Router			/tags [post]
//
// タグを新規作成
func (h *TagHandler) Create(ctx echo.Context) error {
	// ユーザーID取得
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return err
	}
	// リクエスト取得
	request := TagCreateRequest{}
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	err = h.tagCmdService.CreateTag(userId, service.TagCreateCommandInput{
		Name: request.Name,
	})
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
