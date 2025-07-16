package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

const (
	USER_QUERY_LIMIT = 100
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
	Role      string
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
		Role:      user.Role().String(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}

type UserAllQueryInput struct {
	Limit uint
	Page  uint
}

// ユーザーを一覧取得
func (s *UserQueryService) GetUsers(userId uuid.UUID, input UserAllQueryInput) ([]UserQueryOutput, error) {
	users, err := s.userRepository.FindAll(userId, input.Limit, input.Page)
	if err != nil {
		return nil, err
	}
	outputList := []UserQueryOutput{}
	for _, user := range users {
		outputList = append(outputList, UserQueryOutput{
			Id:        user.Id(),
			Name:      user.Name(),
			Email:     user.Email().String(),
			Role:      user.Role().String(),
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
		})
	}
	return outputList, nil
}
