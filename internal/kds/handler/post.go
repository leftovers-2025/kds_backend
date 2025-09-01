package handler

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/service"
)

var (
	ErrPostGetInvalidLimit = common.NewValidationError(errors.New("invalid parameter 'limit'"))
	ErrPostGetInvalidPage  = common.NewValidationError(errors.New("invalid parameter 'page'"))
	ErrPostGetInvalidParam = common.NewValidationError(errors.New("invalid parameter"))
)

type PostHandler struct {
	postCmdService   *service.PostCommandService
	postQueryService *service.PostQueryService
}

func NewPostHandler(
	postCmdService *service.PostCommandService,
	postQueryService *service.PostQueryService,
) *PostHandler {
	if postCmdService == nil {
		panic("nil PostCommandService")
	}
	if postQueryService == nil {
		panic("nil PostQueryService")
	}
	return &PostHandler{
		postCmdService:   postCmdService,
		postQueryService: postQueryService,
	}
}

type PostResponse struct {
	Id          string                   `json:"id"`
	UserId      string                   `json:"userId"`
	Description string                   `json:"description"`
	Location    PostResponseLocationItem `json:"location"`
	Tags        []PostResponseTagItem    `json:"tags"`
	Images      []string                 `json:"images"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
}

type PostResponseTagItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PostResponseLocationItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Get godoc
//
//	@Summary		Get posts
//	@Description	Get a list of posts with optional filters
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			limit		query	int		true	"Limit"
//	@Param			page		query	int		true	"Page"
//	@Param			query		query	string	false	"Query word"
//	@Param			tag			query	string	false	"Tag"
//	@Param			location	query	string	false	"Location"
//	@Param			order		query	string	false	"Order"
//	@Param			orderAsc	query	bool	false	"Order ascending"
//	@Success		200			{array}	PostResponse
//	@Router			/posts [get]
//
// 投稿を取得
func (h *PostHandler) Get(ctx echo.Context) error {
	// パラメーター取得
	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		return ErrPostGetInvalidLimit
	}
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return ErrPostGetInvalidPage
	}
	queryWord := ctx.QueryParam("query")
	tag := ctx.QueryParam("tag")
	location := ctx.QueryParam("location")
	order := ctx.QueryParam("order")
	orderAsc := false
	if ctx.QueryParams().Has("orderAsc") {
		orderAsc, err = strconv.ParseBool(ctx.QueryParam("orderAsc"))
		if err != nil {
			return ErrPostGetInvalidParam
		}
	}
	// 投稿一覧取得
	outputList, err := h.postQueryService.GetPosts(service.PostQueryInput{
		QueryWord: queryWord,
		Tag:       tag,
		Location:  location,
		Order:     order,
		OrderAsc:  orderAsc,
		Limit:     uint(limit),
		Page:      uint(page),
	})
	if err != nil {
		return err
	}
	// レスポンスにマッピング
	responseList := []PostResponse{}
	for _, outputItem := range outputList {
		tags := []PostResponseTagItem{}
		for _, tag := range outputItem.Tags {
			tags = append(tags, PostResponseTagItem{
				Id:   tag.Id.String(),
				Name: tag.Name,
			})
		}
		responseList = append(responseList, PostResponse{
			Id:          outputItem.Id.String(),
			UserId:      outputItem.UserId.String(),
			Description: outputItem.Description,
			Location: PostResponseLocationItem{
				Id:   outputItem.Location.Id.String(),
				Name: outputItem.Location.Name,
			},
			Tags:      tags,
			Images:    outputItem.Images,
			CreatedAt: outputItem.CreatedAt,
			UpdatedAt: outputItem.UpdatedAt,
		})
	}
	return ctx.JSON(http.StatusOK, &responseList)
}

// Create godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post
//	@Tags			posts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			description	formData	string	true	"Description"
//	@Param			locationId	formData	string	true	"Location ID"
//	@Param			tagIds		formData	string	true	"Tag IDs (comma separated)"
//	@Param			image1		formData	file	false	"Image 1"
//	@Param			image2		formData	file	false	"Image 2"
//	@Param			image3		formData	file	false	"Image 3"
//	@Success		204
//	@Router			/posts [post]
//
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
