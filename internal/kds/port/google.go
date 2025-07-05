package port

import "github.com/leftovers-2025/kds_backend/internal/kds/entity"

type GoogleRepository interface {
	CodeAuthorization(code string) (*entity.GoogleUser, error)
}
