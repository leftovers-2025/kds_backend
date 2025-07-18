package mysql

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type MySqlNotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) port.NotificationRepository {
	return &MySqlNotificationRepository{
		db: db,
	}
}

// 通知対象ユーザーを一覧取得
func (r *MySqlNotificationRepository) FindNotifyUsers(locationId uuid.UUID, tagIds []uuid.UUID) ([]entity.User, error) {
	sql := `
		SELECT 
			users.id, users.name, users.email, users.created_at, users.updated_at
		FROM 
			users
		JOIN notifications
			ON notifications.user_id = users.id AND notifications.enabled = 1
		JOIN tag_notifications
			ON tag_notifications.user_id = users.id
		JOIN location_notifications
			ON location_notifications.user_id = users.id
		JOIN google_ids
			ON google_ids.user_id = users.id
		JOIN roles
			ON roles.user_id = users.id
		WHERE
			 location_notifications.location_id = ? OR tag_notifications.tag_id IN (?)
	`
	tagIdBytesList := [][]byte{}
	for _, tagId := range tagIds {
		tagIdBytesList = append(tagIdBytesList, tagId[:])
	}
	models := []UserAndGoogleIdAndRoleModel{}
	query, args, err := sqlx.In(sql, tagIdBytesList, locationId[:])
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&models, query, args...)
	if err != nil {
		return nil, err
	}
	users := []entity.User{}
	for _, model := range models {
		user, err := modelToUser(&model)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

// 通知を保存
func (r *MySqlNotificationRepository) Save(notification *entity.Notification) error {
	return runInTx(r.db, func(tx *sqlx.Tx) error {
		sql := `
			INSERT INTO notifications(user_id, enabled)
			VALUES (:userId, :enabled)
			ON DUPLICATE KEY UPDATE
				enabled = :enabled
		`
		userId := notification.UserId()
		_, err := tx.NamedExec(sql, map[string]any{
			"userId":  userId[:],
			"enabled": notification.IsEnabled(),
		})
		if err != nil {
			return err
		}
		sql = `
			DELETE FROM tag_notifications
			WHERE user_id = :userId
		`
		_, err = tx.NamedExec(sql, map[string]any{
			"userId": userId[:],
		})
		if err != nil {
			return err
		}
		err = r.insertTagsInTx(tx, userId, notification.Tags())
		if err != nil {
			return err
		}
		sql = `
			DELETE FROM location_notifications
			WHERE user_id = :userId
		`
		_, err = tx.NamedExec(sql, map[string]any{
			"userId": userId[:],
		})
		if err != nil {
			return err
		}
		err = r.insertLocationsInTx(tx, userId, notification.Locations())
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MySqlNotificationRepository) insertTagsInTx(tx *sqlx.Tx, userId uuid.UUID, tags []entity.Tag) error {
	query := `
			INSERT INTO tag_notifications(user_id, tag_id)
				SELECT 
					:userId AS user_id, id AS tag_id
				FROM tags
				WHERE tags.id IN (:tagIds)
		`
	tagIds := [][]byte{}
	for _, tag := range tags {
		tagId := tag.Id()
		tagIds = append(tagIds, tagId[:])
	}
	arg := map[string]any{
		"userId": userId[:],
		"tagIds": tagIds,
	}
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, args...)
	return err
}

func (r *MySqlNotificationRepository) insertLocationsInTx(tx *sqlx.Tx, userId uuid.UUID, locations []entity.Location) error {
	query := `
			INSERT INTO location_notifications(user_id, location_id)
				SELECT 
					:userId AS user_id, id AS location_id
				FROM locations
				WHERE locations.id IN (:locationIds)
		`
	locationIds := [][]byte{}
	for _, location := range locations {
		locationId := location.Id()
		locationIds = append(locationIds, locationId[:])
	}
	arg := map[string]any{
		"userId":      userId[:],
		"locationIds": locationIds,
	}
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, args...)
	return err
}
