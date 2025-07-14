package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

var (
	ErrLocationCreateInvalidPermission = common.NewValidationError(errors.New("invalid permission"))
)

type LocationCommandService struct {
	locationRepository port.LocationRepository
}

func NewLocationCommandService(
	locationRepository port.LocationRepository,
) *LocationCommandService {
	if locationRepository == nil {
		panic("nil LocationRepository")
	}
	return &LocationCommandService{
		locationRepository: locationRepository,
	}
}

type LocationCreateCommandInput struct {
	Name string
}

// ロケーション新規作成
func (s *LocationCommandService) CreateLocation(userId uuid.UUID, input LocationCreateCommandInput) error {
	err := s.locationRepository.Create(userId, func(user *entity.User) (*entity.Location, error) {
		// 権限確認
		if user.Role() == entity.ROLE_STUDENT {
			return nil, ErrLocationCreateInvalidPermission
		}
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		// ロケーション作成
		location, err := entity.NewLocation(id, input.Name)
		if err != nil {
			return nil, common.NewValidationError(err)
		}
		return location, nil
	})
	return err
}
