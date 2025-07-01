package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type UserCommandService struct {
	userRepository port.UserRepository
}

func NewUserCommandService(userRepository port.UserRepository) *UserCommandService {
	if userRepository == nil {
		panic("nil UserRepository")
	}
	return &UserCommandService{
		userRepository: userRepository,
	}
}

type UserCreateCommandInput struct {
	GoogleId string
	Name     string
	Email    string
}

type UserCreateCommandOutput struct {
	Id    uuid.UUID
	Name  string
	Email string
}

// ユーザーを新規作成する
func (s *UserCommandService) CreateUser(cmd UserCreateCommandInput) (*UserCreateCommandOutput, error) {
	// id生成
	id, err := uuid.NewV7()
	if err != nil {
		return nil, common.NewInternalError(err)
	}
	// ユーザー作成
	user, err := entity.NewUser(
		id,
		cmd.Name,
		cmd.Email,
		cmd.GoogleId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, common.NewInternalError(err)
	}
	// リポジトリで作成
	err = s.userRepository.Create(user)
	if err != nil {
		return nil, common.NewInternalError(err)
	}
	return &UserCreateCommandOutput{
		Id:    user.Id(),
		Name:  user.Name(),
		Email: user.Email().String(),
	}, nil
}
