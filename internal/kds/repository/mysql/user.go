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

type UserAndGoogleIdAndRoleModel struct {
	Id        []byte    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	GoogleId  string    `db:"google_id"`
	Role      string    `db:"role"`
}

// ユーザーを新規作成する
func (r *MySqlUserRepository) Create(user *entity.User) error {
	return RunInTx(r.db, func(tx *sqlx.Tx) error {
		// ユーザー挿入
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
		// GoogleId挿入
		sql = `
			INSERT INTO google_ids(user_id, google_id)
			VALUES(:userId, :googleId)
		`
		_, err = tx.NamedExec(sql, map[string]any{
			"userId":   userId[:],
			"googleId": user.GoogleId(),
		})
		if err != nil {
			return err
		}
		// ロール挿入
		sql = `
			INSERT INTO roles(user_id, role, updated_at)
			VALUES(:userId, :role, :updatedAt)
		`
		_, err = tx.NamedExec(sql, map[string]any{
			"userId":    userId[:],
			"role":      user.Role().String(),
			"updatedAt": user.UpdatedAt(),
		})
		return err
	})
}

// 対象ユーザの情報を編集する
func (r *MySqlUserRepository) EditUser(userId, targetUserId uuid.UUID, editFn func(user, targetUserId *entity.User) error) error {
	err := RunInTx(r.db, func(tx *sqlx.Tx) error {
		// ユーザ取得
		user, err := getUserInTx(tx, userId)
		if err != nil {
			return err
		}
		// 編集対象ユーザー取得
		targetUser, err := getUserInTx(tx, targetUserId)
		if err != nil {
			return err
		}
		// 編集
		err = editFn(user, targetUser)
		if err != nil {
			return err
		}
		// ユーザーロール情報更新
		query := `
			UPDATE roles
			SET 
				role = :role,
				updated_at = :updatedAt
			WHERE 
				user_id = :userId
		`
		_, err = tx.NamedExec(query, map[string]any{
			"userId":    targetUserId[:],
			"role":      targetUser.Role().String(),
			"updatedAt": targetUser.UpdatedAt(),
		})
		return err
	})
	return err
}

// Idからユーザーを検索する
func (r *MySqlUserRepository) FindById(id uuid.UUID) (*entity.User, error) {
	sql := `
		SELECT 
			users.id, users.name, users.email, users.created_at, users.updated_at, google_ids.google_id, roles.role
		FROM users
			JOIN google_ids
				ON users.id = google_ids.user_id
			JOIN roles
				ON users.id = roles.user_id
		WHERE
			users.id = :id
	`
	model := UserAndGoogleIdAndRoleModel{}
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
	err = row.StructScan(&model)
	if err != nil {
		return nil, err
	}
	// モデルをユーザーに変換
	return modelToUser(&model)
}

// GoogleIdからユーザーを検索する
func (r *MySqlUserRepository) FindByGoogleId(googleId string) (*entity.User, error) {
	sql := `
		SELECT 
			users.id, users.name, users.email, users.created_at, users.updated_at, google_ids.google_id, roles.role
		FROM users
			JOIN google_ids
				ON users.id = google_ids.user_id
			JOIN roles
				ON users.id = roles.user_id
		WHERE
			google_ids.google_id = :googleId
	`
	model := UserAndGoogleIdAndRoleModel{}
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
	return modelToUser(&model)
}

// ユーザーを一覧取得
func (r *MySqlUserRepository) FindAll(userId uuid.UUID, limit, page uint) ([]entity.User, error) {
	users := []entity.User{}
	err := RunInTx(r.db, func(tx *sqlx.Tx) error {
		models := []UserAndGoogleIdAndRoleModel{}
		sql := `
			SELECT 
				users.id, users.name, users.email, users.created_at, users.updated_at, google_ids.google_id, roles.role
			FROM users
				JOIN google_ids
					ON users.id = google_ids.user_id
				JOIN roles
					ON users.id = roles.user_id
			ORDER BY id ASC
			LIMIT ? OFFSET ?
		`
		// クエリ実行ユーザー取得
		user, err := getUserInTx(tx, userId)
		if err != nil {
			return err
		}
		// 閲覧可能か確認
		err = user.CanSeeUsers()
		if err != nil {
			return err
		}
		// 一覧取得
		err = tx.Select(&models, sql, limit, (page-1)*limit)
		if err != nil {
			return err
		}
		// ユーザーにマッピング
		for _, model := range models {
			user, err := modelToUser(&model)
			if err != nil {
				return err
			}
			users = append(users, *user)
		}
		return nil
	})
	return users, err
}

// トランザクションでユーザーを取得
func getUserInTx(tx *sqlx.Tx, userId uuid.UUID) (*entity.User, error) {
	query := `
		SELECT 
			users.id, users.email, users.name, users.created_at, users.updated_at,
			google_ids.google_id,
			roles.role
		FROM users
			JOIN google_ids
				ON google_ids.user_id = users.id
			JOIN roles
				ON roles.user_id = users.id
		WHERE
			users.id = ?
	`
	model := UserAndGoogleIdAndRoleModel{}
	err := tx.Get(&model, query, userId[:])
	if err != nil {
		return nil, err
	}
	user, err := modelToUser(&model)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UserAndGoogleIdDBModelをUserに変換する
func modelToUser(model *UserAndGoogleIdAndRoleModel) (*entity.User, error) {
	id, err := uuid.FromBytes(model.Id)
	if err != nil {
		return nil, err
	}
	email, err := entity.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}
	role := entity.RoleFromString(model.Role)
	return entity.NewUser(
		id,
		model.Name,
		model.GoogleId,
		email,
		role,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
