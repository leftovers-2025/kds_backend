package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

type LocationHandler struct {
	locationCmdService   *service.LocationCommandService
	locationQueryService *service.LocationQueryService
}

func NewLocationHandler(
	locationCmdService *service.LocationCommandService,
	locationQueryService *service.LocationQueryService,
) *LocationHandler {
	if locationCmdService == nil {
		panic("nil LocationCommandService")
	}
	if locationQueryService == nil {
		panic("nil LocationQueryService")
	}
	return &LocationHandler{
		locationCmdService:   locationCmdService,
		locationQueryService: locationQueryService,
	}
}

// ロケーションレスポンス
type LocationResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// GetAll godoc
//
//	@Summary		Get all locations
//	@Description	Get a list of all locations
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	LocationResponse
//	@Router			/locations [get]
//
// タグを一覧取得
func (h *LocationHandler) GetAll(ctx echo.Context) error {
	locations, err := h.locationQueryService.GetAllLocations()
	if err != nil {
		return err
	}
	responseList := []LocationResponse{}
	for _, location := range locations {
		responseList = append(responseList, LocationResponse{
			Id:   location.Id.String(),
			Name: location.Name,
		})
	}
	return ctx.JSON(http.StatusOK, &responseList)
}

// ロケーション作成リクエスト
type LocationCreateRequest struct {
	Name string `json:"name"`
}

// Create godoc
//
//	@Summary		Create a new location
//	@Description	Create a new location
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	LocationCreateRequest	true	"Location create request"
//	@Success		204
//	@Router			/locations [post]
//
// ロケーションを新規作成
func (h *LocationHandler) Create(ctx echo.Context) error {
	// ユーザーID取得
	userId, err := getUserIdFromCtx(ctx)
	if err != nil {
		return err
	}
	// リクエスト取得
	request := LocationCreateRequest{}
	if err := ctx.Bind(&request); err != nil {
		return common.NewValidationError(err)
	}
	// ロケーション作成
	err = h.locationCmdService.CreateLocation(userId, service.LocationCreateCommandInput{
		Name: request.Name,
	})
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
