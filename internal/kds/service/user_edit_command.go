package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

var (
	ErrUserEditNoPermission      = common.NewValidationError(errors.New("permission error"))
	ErrUserEditInvalidRole       = common.NewValidationError(errors.New("invalid role"))
	ErrUserEditInvalidTargetRole = common.NewValidationError(errors.New("target user`s role is invalid"))
	ErrUserEditInvalidUsers      = errors.New("invalid users")
)

type UserEditCommandService struct {
	userRepository port.UserRepository
}

func NewUserEditCommandService(
	userRepository port.UserRepository,
) *UserEditCommandService {
	if userRepository == nil {
		panic("nil UserRepository")
	}
	return &UserEditCommandService{
		userRepository: userRepository,
	}
}

type UserEditRoleCommandInput struct {
	TargetUserId uuid.UUID
	Role         string
}

type UserTranferRootCommandInput struct {
	TargetUserId uuid.UUID
}

// 対象ユーザーのロールを編集する
func (u *UserEditCommandService) EditRole(userId uuid.UUID, input UserEditRoleCommandInput) error {
	// ロール取得
	newRole := entity.RoleFromString(input.Role)
	if newRole == entity.ROLE_UNKNOWN {
		return ErrUserEditInvalidRole
	}
	// ユーザー編集
	err := u.userRepository.EditUser(userId, input.TargetUserId, func(user, targetUser *entity.User) error {
		// 付与可能な権限か確認
		if !user.Role().CanEdit(newRole) {
			return ErrUserEditNoPermission
		}
		// 対象ユーザーに対して編集できるか確認
		if !user.Role().CanEdit(targetUser.Role()) {
			return ErrUserEditNoPermission
		}
		// 権限更新
		err := targetUser.UpdateRole(newRole)
		return err
	})
	return err
}

// ルート権限を譲渡する
func (u *UserEditCommandService) TransferRoot(userId uuid.UUID, input UserTranferRootCommandInput) error {
	err := u.userRepository.UpdateTwoUsers(userId, input.TargetUserId, func(user1, user2 *entity.User) error {
		// 権限確認
		if !user1.IsRoot() || !user2.IsTeacher() {
			return ErrUserEditNoPermission
		}
		user1.UpdateRole(entity.ROLE_TEACHER)
		user2.UpdateRole(entity.ROLE_ROOT)
		return nil
	})
	return err
}
