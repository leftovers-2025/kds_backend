package port

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

type UserRepository interface {
	Create(*entity.User) error
	FindById(uuid.UUID) (*entity.User, error)
	FindByGoogleId(string) (*entity.User, error)
	FindAll(userId uuid.UUID, limit, page uint) ([]entity.User, error)
	EditUser(userId, targetUserId uuid.UUID, editFn func(user, targetUser *entity.User) error) error
	UpdateTwoUsers(userId, userId2 uuid.UUID, updateFn func(user, user2 *entity.User) error) error
}
