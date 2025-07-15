package port

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

type UserRepository interface {
	Create(*entity.User) error
	FindById(uuid.UUID) (*entity.User, error)
	FindByGoogleId(string) (*entity.User, error)
	EditUser(userId, targetUserId uuid.UUID, editFn func(user, targetUser *entity.User) error) error
}
