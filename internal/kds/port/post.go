package port

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

type PostRepository interface {
	Create(userId, locationId uuid.UUID, tagIds []uuid.UUID, createFn func(*entity.User, *entity.Location, []entity.Tag) (*entity.Post, error)) error
	FindById(uuid.UUID) (*entity.Post, error)
	FindAll(limit, page uint) ([]entity.Post, error)
	FindWithFilter(queryWord, tag, location, order string, orderAsc bool, limit, page uint) ([]entity.Post, error)
}
