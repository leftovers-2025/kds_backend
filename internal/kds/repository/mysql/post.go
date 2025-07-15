package mysql

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
	"github.com/leftovers-2025/kds_backend/internal/kds/repository/s3"
)

type MySqlPostRepository struct {
	db *sqlx.DB
	s3 *s3.S3Repository
}

func NewMySqlPostRepository(db *sqlx.DB, s3Repository *s3.S3Repository) port.PostRepository {
	if db == nil {
		panic("nil MySql Database")
	}
	if s3Repository == nil {
		panic("nil S3Repository")
	}
	return &MySqlPostRepository{
		db: db,
		s3: s3Repository,
	}
}

// 投稿を新規作成
func (r *MySqlPostRepository) Create(userId, locationId uuid.UUID, tagIds []uuid.UUID, createFn func(*entity.User, *entity.Location, []entity.Tag) (*entity.Post, error)) error {
	return RunInTx(r.db, func(tx *sqlx.Tx) error {
		// ユーザー取得
		user, err := getUserInTx(tx, userId)
		if err != nil {
			return err
		}
		// ロケーション取得
		location, err := getLocationInTx(tx, locationId)
		if err != nil {
			return err
		}
		// タグ一覧
		tags, err := getTagsInTx(tx, tagIds)
		if err != nil {
			return err
		}
		// 投稿作成
		post, err := createFn(user, location, tags)
		if err != nil {
			return err
		}
		query := `
			INSERT INTO posts(id, user_id, location_id, description, created_at, updated_at)
			VALUES(:id, :userId, :locationId, :description, :createdAt, :updatedAt)
		`
		postId := post.Id()
		userId := post.UserId()
		locationId := post.Location().Id()
		// 投稿を作成
		_, err = tx.NamedExec(query, map[string]any{
			"id":          postId[:],
			"userId":      userId[:],
			"locationId":  locationId[:],
			"description": post.Description(),
			"createdAt":   post.CreatedAt(),
			"updatedAt":   post.UpdatedAt(),
		})
		if err != nil {
			return err
		}
		// タグを作成
		if len(post.Tags()) > 0 {
			query = `
				INSERT INTO post_tags(post_id, tag_id)
				VALUES (:postId, :tagId)
			`
			tagMap := []map[string]any{}
			// 全てのタグをマップに追加
			for _, tag := range post.Tags() {
				tagId := tag.Id()
				tagMap = append(tagMap, map[string]any{
					"postId": postId[:],
					"tagId":  tagId[:],
				})
			}
			// タグ追加SQL実行
			_, err = tx.NamedExec(query, tagMap)
			if err != nil {
				return err
			}
		}
		// 画像を作成
		if len(post.Images()) > 0 {
			query = `
				INSERT INTO post_images(post_id, image_url)
				VALUES (:postId, :imageUrl)
			`
			imageMap := []map[string]any{}
			imageModels := []ImageModel{}

			// 全てのイメージをマップに追加
			for _, image := range post.Images() {
				// ファイルを取得
				file, err := image.File()
				if err != nil {
					return err
				}
				// モデル作成
				imageModel := ImageModel{
					Name: uuid.NewString() + ".jpg",
					File: file,
				}
				// モデル一覧に追加
				imageModels = append(imageModels, imageModel)
				// マップに追加
				imageMap = append(imageMap, map[string]any{
					"postId":   postId[:],
					"imageUrl": imageModel.Name,
				})
			}
			// イメージ追加SQL実行
			_, err = tx.NamedExec(query, imageMap)
			if err != nil {
				return err
			}
			// 画像をアップロード
			err = r.uploadImages(imageModels)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Idから投稿を取得
func (r *MySqlPostRepository) FindById(postId uuid.UUID) (*entity.Post, error) {
	return nil, nil
}

type ImageModel struct {
	Name string
	File *multipart.FileHeader
}

// イメージアップロード
func (r *MySqlPostRepository) uploadImages(models []ImageModel) error {
	for _, imageModel := range models {
		err := r.s3.UploadImage(imageModel.Name, imageModel.File)
		if err != nil {
			return err
		}
	}
	return nil
}
