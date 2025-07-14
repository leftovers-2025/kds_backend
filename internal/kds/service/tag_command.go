package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

var (
	ErrTagCreateInvalidPermission = common.NewValidationError(errors.New("invalid permission"))
)

type TagCommandService struct {
	tagRepository port.TagRepository
}

func NewTagCommandService(tagRepository port.TagRepository) *TagCommandService {
	if tagRepository == nil {
		panic("nil TagRepository")
	}
	return &TagCommandService{
		tagRepository: tagRepository,
	}
}

type TagCreateCommandInput struct {
	Name string
}

// タグ新規作成
func (s *TagCommandService) CreateTag(userId uuid.UUID, input TagCreateCommandInput) error {
	// リポジトリ保存
	return s.tagRepository.Create(userId, func(user *entity.User) (*entity.Tag, error) {
		// 権限確認
		if user.Role() == entity.ROLE_STUDENT {
			return nil, ErrTagCreateInvalidPermission
		}
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		// タグ作成
		tag, err := entity.NewTag(id, input.Name)
		if err != nil {
			return nil, common.NewValidationError(err)
		}
		return tag, nil
	})
}
