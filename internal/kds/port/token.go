package port

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

type TokenRepository interface {
	AddWhitelist(*entity.RefreshToken) error
	RemoveWhitelist(uuid.UUID) error
	InWhitelist(uuid.UUID) (bool, error)
}
