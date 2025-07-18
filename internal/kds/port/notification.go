package port

import (
	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

type NotificationRepository interface {
	FindNotifyUsers(locationId uuid.UUID, tagIds []uuid.UUID) ([]entity.User, error)
	Save(notification *entity.Notification) error
}
