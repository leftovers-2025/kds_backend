package port

import "github.com/leftovers-2025/kds_backend/internal/kds/entity"

type EmailRepository interface {
	SendAll(users []entity.User, subject, content string) error
}
