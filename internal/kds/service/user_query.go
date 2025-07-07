package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type UserQueryService struct {
	userRepository port.UserRepository
}

func NewUserQueryService(userRepository port.UserRepository) *UserQueryService {
	if userRepository == nil {
		panic("nil UserRepository")
	}
	return &UserQueryService{
		userRepository: userRepository,
	}
}

type UserQueryOutput struct {
	Id        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Idからユーザー情報を取得
func (s *UserQueryService) GetUser(userId uuid.UUID) (*UserQueryOutput, error) {
	user, err := s.userRepository.FindById(userId)
	if err != nil {
		return nil, err
	}
	return &UserQueryOutput{
		Id:        user.Id(),
		Name:      user.Name(),
		Email:     user.Email().String(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}
