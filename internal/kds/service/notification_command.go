package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

var (
	ErrNotiicationTagNotFound      = common.NewValidationError(errors.New("notification tags not found"))
	ErrNotiicationLocationNotFound = common.NewValidationError(errors.New("notification locations not found"))
)

type NotificationCommandService struct {
	tagRepository         port.TagRepository
	emailRepository       port.EmailRepository
	locationRepository    port.LocationRepository
	notificationRepsitory port.NotificationRepository
}

func NewNotificationCommandService(
	tagRepository port.TagRepository,
	emailRepository port.EmailRepository,
	locationRepository port.LocationRepository,
	notificationRepsitory port.NotificationRepository,
) *NotificationCommandService {
	if tagRepository == nil {
		panic("nil TagRepository")
	}
	if emailRepository == nil {
		panic("nil EmailRepository")
	}
	if locationRepository == nil {
		panic("nil LocationRepository")
	}
	if notificationRepsitory == nil {
		panic("nil NotificationRepository")
	}
	return &NotificationCommandService{
		tagRepository:         tagRepository,
		emailRepository:       emailRepository,
		locationRepository:    locationRepository,
		notificationRepsitory: notificationRepsitory,
	}

}

type NotifyCommandInput struct {
	PostId     uuid.UUID
	LocationId uuid.UUID
	TagIds     []uuid.UUID
}

// 対象のユーザーに落とし物通知を送信
func (s *NotificationCommandService) Notify(input NotifyCommandInput) error {
	users, err := s.notificationRepsitory.FindNotifyUsers(input.LocationId, input.TagIds)
	if err != nil {
		return err
	}
	err = s.emailRepository.SendAll(users, "類似する落とし物が見つかりました", "http://localhost:8630/posts/"+input.PostId.String())
	return err
}

type NotificationCommandInput struct {
	UserId      uuid.UUID
	Enabled     bool
	LocationIds []uuid.UUID
	TagIds      []uuid.UUID
}

// 通知設定を保存する
func (s *NotificationCommandService) SaveNotification(input NotificationCommandInput) error {
	var err error
	// タグ取得
	tags := []entity.Tag{}
	if len(input.TagIds) > 0 {
		tags, err = s.tagRepository.FindByIds(input.TagIds)
		if err != nil {
			return err
		}
		if len(tags) != len(input.TagIds) {
			return ErrNotiicationTagNotFound
		}
	}
	// ロケーション取得
	locations := []entity.Location{}
	if len(input.LocationIds) > 0 {
		locations, err = s.locationRepository.FindByIds(input.LocationIds)
		if err != nil {
			return err
		}
		if len(locations) != len(input.LocationIds) {
			return ErrNotiicationLocationNotFound
		}
	}
	notification, err := entity.NewNotificaton(input.UserId, input.Enabled, locations, tags)
	if err != nil {
		return err
	}
	return s.notificationRepsitory.Save(notification)
}
