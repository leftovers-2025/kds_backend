package port

import "github.com/leftovers-2025/kds_backend/internal/kds/entity"

type UserRepository interface {
	Create(*entity.User) error
}
