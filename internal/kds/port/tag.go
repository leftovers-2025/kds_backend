package port

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

type TagRepository interface {
	Create(userId uuid.UUID, createFn func(*entity.User) (*entity.Tag, error)) error
	FindAll() ([]entity.Tag, error)
}
