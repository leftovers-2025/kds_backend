package service

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

var (
	ErrPostCreateInvalidPermission = common.NewValidationError(errors.New("permission error"))
	ErrPostCreateTagNotFound       = common.NewValidationError(errors.New("tag not found"))
	ErrPostCreateLocationNotFound  = common.NewValidationError(errors.New("location not found"))
)

type PostCommandService struct {
	postRepository port.PostRepository
}

func NewPostCommandService(
	postRepository port.PostRepository,
) *PostCommandService {
	if postRepository == nil {
		panic("nil PostRepository")
	}
	return &PostCommandService{
		postRepository: postRepository,
	}
}

type PostCreateCommandInput struct {
	LocationId  uuid.UUID
	Description string
	TagIds      []uuid.UUID
	Images      []multipart.FileHeader
}

// 投稿を新規作成する
func (s *PostCommandService) CreatePost(userId uuid.UUID, input PostCreateCommandInput) error {
	// リポジトリ保存
	err := s.postRepository.Create(
		userId,
		input.LocationId,
		input.TagIds,
		func(user *entity.User, location *entity.Location, tags []entity.Tag) (*entity.Post, error) {
			// 権限確認
			if user.Role() == entity.ROLE_STUDENT {
				return nil, ErrPostCreateInvalidPermission
			}
			// ロケーション存在確認
			if location == &entity.NilLocation {
				return nil, ErrPostCreateLocationNotFound
			}
			// タグ一致確認
			if len(tags) != len(input.TagIds) {
				return nil, ErrPostCreateTagNotFound
			}
			// 画像作成
			images := []entity.Image{}
			for _, imageFile := range input.Images {
				image, err := entity.NewFileImage(&imageFile)
				if err != nil {
					return nil, err
				}
				images = append(images, *image)
			}
			// id生成
			id, err := uuid.NewV7()
			if err != nil {
				return nil, err
			}
			// 投稿作成
			post, err := entity.NewPost(
				id,
				user.Id(),
				*location,
				input.Description,
				tags,
				images,
				time.Now(),
				time.Now(),
			)
			if err != nil {
				return nil, err
			}
			return post, nil
		})
	return err
}
