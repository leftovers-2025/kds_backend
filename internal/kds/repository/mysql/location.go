package mysql

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type MySqlLocationRepository struct {
	db *sqlx.DB
}

func NewMySqlLocationRepository(db *sqlx.DB) port.LocationRepository {
	return &MySqlLocationRepository{
		db: db,
	}
}

type LocationModel struct {
	Id   []byte `db:"id"`
	Name string `db:"name"`
}

// ロケーションを新規作成
func (r *MySqlLocationRepository) Create(userId uuid.UUID, createFn func(*entity.User) (*entity.Location, error)) error {
	return RunInTx(r.db, func(tx *sqlx.Tx) error {
		// ユーザー取得
		user, err := getUserInTx(tx, userId)
		if err != nil {
			return err
		}
		// ロケーション作成
		location, err := createFn(user)
		if err != nil {
			return err
		}
		query := `
			INSERT IGNORE INTO locations(id, name)
			VALUES (:id, :name)
		`
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		// DB挿入
		_, err = r.db.NamedExec(query, map[string]any{
			"id":   id[:],
			"name": location.Name(),
		})
		return err
	})
}

// ロケーションを一覧取得
func (r *MySqlLocationRepository) FindAll() ([]entity.Location, error) {
	query := `
		SELECT 
			id, name
		FROM locations
	`
	models := []LocationModel{}
	err := r.db.Select(&models, query)
	if err != nil {
		return nil, err
	}
	locations := []entity.Location{}
	for _, model := range models {
		id, err := uuid.FromBytes(model.Id)
		if err != nil {
			return nil, err
		}
		location, err := entity.NewLocation(id, model.Name)
		if err != nil {
			return nil, err
		}
		locations = append(locations, *location)
	}
	return locations, nil
}

// トランザクションでロケーションを取得
func getLocationInTx(tx *sqlx.Tx, locationid uuid.UUID) (*entity.Location, error) {
	query := `
		SELECT 
			id, name
		FROM locations
		WHERE 
			id = ?
	`
	model := LocationModel{}
	err := tx.Get(&model, query, locationid[:])
	if err != nil {
		return nil, err
	}
	id, err := uuid.FromBytes(model.Id)
	if err != nil {
		return nil, err
	}
	location, err := entity.NewLocation(id, model.Name)
	if err != nil {
		return nil, err
	}
	return location, nil
}
