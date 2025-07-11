package port

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

type LocationRepository interface {
	Create(userId uuid.UUID, createFn func(*entity.User) (*entity.Location, error)) error
	FindAll() ([]entity.Location, error)
}
