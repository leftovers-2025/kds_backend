package service

import (
	"time"

	"github.com/google/uuid"
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
	email, err := entity.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}
	// id生成
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	// ロール設定
	userRole := entity.ROLE_STUDENT
	if email.IsTeacher() {
		userRole = entity.ROLE_TEACHER
	}
	// ユーザー作成
	user, err := entity.NewUser(
		id,
		cmd.Name,
		cmd.GoogleId,
		email,
		userRole,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	// リポジトリで作成
	err = s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	return &UserCreateCommandOutput{
		Id:    user.Id(),
		Name:  user.Name(),
		Email: user.Email().String(),
	}, nil
}
