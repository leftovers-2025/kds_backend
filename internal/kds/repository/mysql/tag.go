package mysql

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

type MySqlTagRepository struct {
	db *sqlx.DB
}

func NewMySqlTagRepository(db *sqlx.DB) port.TagRepository {
	return &MySqlTagRepository{
		db: db,
	}
}

type TagModel struct {
	Id   []byte `db:"id"`
	Name string `db:"name"`
}

// タグを新規作成
func (r *MySqlTagRepository) Create(userId uuid.UUID, createFn func(user *entity.User) (*entity.Tag, error)) error {
	err := RunInTx(r.db, func(tx *sqlx.Tx) error {
		// ユーザー作成
		user, err := getUserInTx(tx, userId)
		if err != nil {
			return err
		}
		// タグ作成
		tag, err := createFn(user)
		if err != nil {
			return err
		}
		query := `
			INSERT IGNORE INTO tags(id, name)
			VALUES (:id, :name)
		`
		tagId := tag.Id()
		// DB挿入
		_, err = r.db.NamedExec(query, map[string]any{
			"id":   tagId[:],
			"name": tag.Name(),
		})
		return err
	})
	return err
}

// タグを一覧取得
func (r *MySqlTagRepository) FindAll() ([]entity.Tag, error) {
	query := `
		SELECT 
			id, name
		FROM tags
	`
	tags := []TagModel{}
	// クエリ発行
	err := r.db.Select(&tags, query)
	if err != nil {
		return nil, err
	}
	return tagModelsToEntities(tags)
}

// トランザクションでタグを一覧取得
func getTagsInTx(tx *sqlx.Tx, tagIds []uuid.UUID) ([]entity.Tag, error) {
	sql := `
		SELECT 
			id, name
		FROM tags
		WHERE 
			tags.id IN (?)
	`
	ids := [][]byte{}
	for _, tagId := range tagIds {
		ids = append(ids, tagId[:])
	}
	tags := []TagModel{}
	// INクエリを組み立て
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return nil, err
	}
	// クエリ発行
	err = tx.Select(&tags, query, args...)
	if err != nil {
		return nil, err
	}
	return tagModelsToEntities(tags)
}

// タグモデルのリストをエンティティに変換
func tagModelsToEntities(tagModels []TagModel) ([]entity.Tag, error) {
	tags := []entity.Tag{}
	for _, model := range tagModels {
		id, err := uuid.FromBytes(model.Id)
		if err != nil {
			return nil, err
		}
		tag, err := entity.NewTag(id, model.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}
