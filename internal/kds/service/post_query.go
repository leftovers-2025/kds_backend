package service

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

const (
	POST_QUERY_LIMIT            = 50
	POST_QUERT_ORDER_ASC        = false
	POST_QUERY_QUERY_WORD_LIMIT = 32
)

var (
	ErrPostQueryInvalidOrder = common.NewValidationError(errors.New("invalid query order"))
	ErrPostQueryInvalidLimit = common.NewValidationError(errors.New("invalid query limit"))
	ErrPostQueryInvalidPage  = common.NewValidationError(errors.New("invalid query page"))
	queryOrders              = []string{
		"createdAt", "location", "userId",
	}
)

type PostQueryService struct {
	postRepository port.PostRepository
}

func NewPostQueryService(postRepositorty port.PostRepository) *PostQueryService {
	if postRepositorty == nil {
		panic("nil PostRepository")
	}
	return &PostQueryService{
		postRepository: postRepositorty,
	}
}

type PostQueryInput struct {
	QueryWord string
	Tag       string
	Location  string
	Order     string
	OrderAsc  bool
	Limit     uint
	Page      uint
}

type PostQueryOutput struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	Description string
	Location    PostQueryOutputLocation
	Tags        []PostQueryOutputTag
	Images      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostQueryOutputLocation struct {
	Id   uuid.UUID
	Name string
}

type PostQueryOutputTag struct {
	Id   uuid.UUID
	Name string
}

// 投稿を一覧取得
func (s *PostQueryService) GetPosts(input PostQueryInput) ([]PostQueryOutput, error) {
	if input.Limit == 0 || input.Limit > POST_QUERY_LIMIT {
		return nil, ErrPostQueryInvalidLimit
	}
	if input.Page == 0 {
		return nil, ErrPostQueryInvalidPage
	}
	if err := s.isValidQueryOrder(input.Order); err != nil {
		return nil, err
	}
	// 投稿を検索
	posts, err := s.postRepository.FindWithFilter(input.QueryWord, input.Tag, input.Location, input.Order, input.OrderAsc, input.Limit, input.Page)
	if err != nil {
		return nil, err
	}
	outputList := []PostQueryOutput{}
	// 出力型に変換
	for _, post := range posts {
		// タグを変換
		tags := []PostQueryOutputTag{}
		for _, tag := range post.Tags() {
			tags = append(tags, PostQueryOutputTag{
				Id:   tag.Id(),
				Name: tag.Name(),
			})
		}
		// 画像を変換
		images := []string{}
		for _, image := range post.Images() {
			name, _ := image.Name()
			images = append(images, name)
		}
		// リストに追加
		outputList = append(outputList, PostQueryOutput{
			Id:          post.Id(),
			UserId:      post.UserId(),
			Description: post.Description(),
			Location: PostQueryOutputLocation{
				Id:   post.Location().Id(),
				Name: post.Location().Name(),
			},
			Tags:      tags,
			Images:    images,
			CreatedAt: post.CreatedAt(),
			UpdatedAt: post.UpdatedAt(),
		})
	}
	return outputList, nil
}

// ソート対象が正しいかを確認
func (s *PostQueryService) isValidQueryOrder(order string) error {
	if order == "" {
		return nil
	}
	if slices.Contains(queryOrders, order) {
		return nil
	}
	return ErrPostQueryInvalidOrder
}
