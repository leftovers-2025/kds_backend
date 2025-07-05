package mysql

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type MySqlUserRepository struct {
	db *sqlx.DB
}

func NewMySqlUserRepository(db *sqlx.DB) port.UserRepository {
	if db == nil {
		panic("nil MySQL DB")
	}
	return &MySqlUserRepository{
		db: db,
	}
}

type UserAndGoogleIdDBModel struct {
	Id        []byte    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	GoogleId  string    `db:"google_id"`
}

// ユーザーを新規作成する
func (r *MySqlUserRepository) Create(user *entity.User) error {
	return RunInTx(r.db, func(tx *sqlx.Tx) error {
		sql := `
			INSERT INTO users(id, name, email, created_at, updated_at)
			VALUES(:id, :name, :email, :createdAt, :updatedAt)
		`
		userId := user.Id()
		_, err := tx.NamedExec(sql, map[string]any{
			"id":        userId[:],
			"name":      user.Name(),
			"email":     user.Email(),
			"createdAt": user.CreatedAt(),
			"updatedAt": user.UpdatedAt(),
		})
		if err != nil {
			return err
		}
		sql = `
			INSERT INTO google_ids(user_id, google_id)
			VALUES(:userId, :googleId)
		`
		_, err = tx.NamedExec(sql, map[string]any{
			"userId":   userId[:],
			"googleId": user.GoogleId(),
		})
		return err
	})
}

// Idからユーザーを検索する
func (r *MySqlUserRepository) FindById(id uuid.UUID) (*entity.User, error) {
	sql := `
		SELECT 
			users.id, users.name, users.email, users.created_at, users.updated_at, google_ids.google_id
		FROM users
			JOIN google_ids
			ON users.id = google_ids.user_id
		WHERE
			users.id = :id
	`
	model := UserAndGoogleIdDBModel{}
	row, err := r.db.NamedQuery(sql, map[string]any{
		"id": id[:],
	})
	if err != nil {
		return nil, err
	}
	// ユーザーが返ってきかた確認
	if !row.Next() {
		return nil, common.NewNotFoundError(fmt.Errorf("user id %s not found", id.String()))
	}
	// モデルにバインド
	columns, err := row.Columns()
	if err != nil {
		fmt.Printf("columns err: %s\n", err.Error())
	} else {
		fmt.Println(columns)
	}
	err = row.StructScan(&model)
	if err != nil {
		return nil, err
	}
	// モデルをユーザーに変換
	return userAndGoogleIdModelToUser(&model)
}

// GoogleIdからユーザーを検索する
func (r *MySqlUserRepository) FindByGoogleId(googleId string) (*entity.User, error) {
	sql := `
		SELECT 
			users.id, users.name, users.email, users.created_at, users.updated_at, google_ids.google_id
		FROM users
			JOIN google_ids
			ON users.id = google_ids.user_id
		WHERE
			google_ids.google_id = :googleId
	`
	model := UserAndGoogleIdDBModel{}
	row, err := r.db.NamedQuery(sql, map[string]any{
		"googleId": googleId,
	})
	if err != nil {
		return nil, err
	}
	// ユーザーが返ってきかた確認
	if !row.Next() {
		return nil, common.NewNotFoundError(fmt.Errorf("user google_id %s not found", googleId))
	}
	// モデルにバインド
	err = row.StructScan(&model)
	if err != nil {
		return nil, err
	}
	// モデルをユーザーに変換
	return userAndGoogleIdModelToUser(&model)
}

// UserAndGoogleIdDBModelをUserに変換する
func userAndGoogleIdModelToUser(model *UserAndGoogleIdDBModel) (*entity.User, error) {
	id, err := uuid.FromBytes(model.Id)
	if err != nil {
		return nil, err
	}
	return entity.NewUser(
		id,
		model.Name,
		model.Email,
		model.GoogleId,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
