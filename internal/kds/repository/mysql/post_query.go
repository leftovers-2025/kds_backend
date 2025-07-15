package mysql

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/leftovers-2025/kds_backend/internal/kds/common"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
)

var (
	ErrPostQueryInvalidOrderName = common.NewValidationError(errors.New("invalid query order name"))
)

type PostModel struct {
	Id           []byte         `db:"id"`
	UserId       []byte         `db:"user_id"`
	LocationId   []byte         `db:"location_id"`
	LocationName string         `db:"location_name"`
	TagId        *[]byte        `db:"tag_id"`
	TagName      sql.NullString `db:"tag_name"`
	Description  string         `db:"description"`
	Image        sql.NullString `db:"image_url"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

// 条件付きで投稿を一覧取得
func (r *MySqlPostRepository) FindWithFilter(queryWord, tag, location, order string, orderAsc bool, limit, page uint) ([]entity.Post, error) {
	if queryWord == "" && tag == "" && location == "" && order == "" && !orderAsc {
		return r.FindAll(limit, page)
	} else {
		return r.findWithFilter(queryWord, tag, location, order, orderAsc, limit, page)
	}
}

// 条件付きで投稿を一覧取得
func (r *MySqlPostRepository) findWithFilter(queryWord, tag, location, order string, orderAsc bool, limit, page uint) ([]entity.Post, error) {
	// ソート対象のカラム名を取得
	orderColumn, err := r.toOrderColumn(order)
	if err != nil {
		return nil, err
	}
	// SQLにバインドするパラメーターリスト
	params := []any{}
	// SQL文
	query := `
		SELECT 
			posts.id, user_id, posts.location_id, location_name, tag_id, tag_name, description, image_url, created_at, updated_at
		FROM 
			(
				SELECT * FROM posts
	`
	// ロケーション検索がある場合
	if location != "" || queryWord != "" {
		// 検索クエリにロケーションを追加
		query += `
				JOIN (
					SELECT locations.id AS l_id, locations.name AS location_name
					FROM locations
				) AS locations
					ON locations.l_id = posts.location_id
		`
	}
	// タグ検索がある場合
	if tag != "" || queryWord != "" {
		// 検索クエリにタグを追加
		query += `LEFT JOIN (
					SELECT post_id, tags.id AS tag_id, name AS tag_name
					FROM post_tags
					JOIN tags
						ON post_tags.tag_id = tags.id
				  ) AS post_tags 
						ON posts.id = post_tags.post_id`
	}
	// 条件絞り込み
	if tag != "" || location != "" || queryWord != "" {
		query += " WHERE"
	}
	// フィルターとワード検索の両方を行う場合
	if (tag != "" || location != "") && queryWord != "" {
		query += " ("
	}
	// タグ絞り込み
	if tag != "" {
		query += " tag_name = ?"
		params = append(params, tag)
	}
	// ロケーション絞り込み
	if location != "" {
		if len(params) > 0 {
			query += " AND"
		}
		query += " location_name = ?"
		params = append(params, location)
	}
	// ワード検索
	if queryWord != "" {
		// フィルター検索していた場合
		if len(params) > 0 {
			query += " ) OR"
		}
		query += `
			posts.description LIKE ?
			OR tag_name LIKE ?
			OR location_name LIKE ?
			`
		params = append(params, "%"+queryWord+"%", queryWord+"%", queryWord+"%")
	}
	// ソートとページネーション
	query += " ORDER BY " + orderColumn
	if orderAsc {
		query += " ASC"
	} else {
		query += " DESC"
	}
	query += ` LIMIT ? OFFSET ?`
	params = append(params, limit, (page-1)*limit)
	// 絞り込み終了
	query += " ) AS posts"
	// ロケーションを絞り込み時に取得していない場合
	if location == "" && queryWord == "" {
		query += `
		JOIN (
			SELECT locations.id AS l_id, locations.name AS location_name
			FROM locations
		) AS locations
		ON locations.l_id = posts.location_id
		`
	}
	// タグを絞り込み時に取得していない場合
	if tag == "" && queryWord == "" {
		query += `
		LEFT JOIN (
			SELECT post_id, tags.id AS tag_id, name AS tag_name
			FROM post_tags
			JOIN tags
			ON post_tags.tag_id = tags.id
		) AS post_tags
			ON posts.id = post_tags.post_id
		`
	}
	// 画像取得クエリ
	query += `
		LEFT JOIN post_images
			ON post_images.post_id = posts.id`
	// ソート
	query += " ORDER BY " + orderColumn
	if orderAsc {
		query += " ASC"
	} else {
		query += " DESC"
	}
	// データ取得
	models := []PostModel{}
	err = r.db.Select(&models, query, params...)
	if err != nil {
		return nil, err
	}
	// マッピングして返す
	return r.modelsToEntities(models)
}

// 投稿を一覧取得
func (r *MySqlPostRepository) FindAll(limit, page uint) ([]entity.Post, error) {
	query := `
		SELECT 
			posts.id, user_id, location_id, locations.name AS location_name, tag_id, post_tags.name AS tag_name, description, image_url, created_at, updated_at
		FROM 
			(
				SELECT * FROM posts
				ORDER BY id DESC
				LIMIT ? OFFSET ?
			) AS posts
		JOIN locations
			ON locations.id = posts.location_id
		LEFT JOIN (
			SELECT post_id, tags.id AS tag_id, name
			FROM post_tags
			JOIN tags
			ON post_tags.tag_id = tags.id
		) AS post_tags
			ON posts.id = post_tags.post_id
		LEFT JOIN post_images
			ON post_images.post_id = posts.id
		ORDER BY id DESC
	`
	models := []PostModel{}
	err := r.db.Select(&models, query, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models)
}

// ソートのカラムに変換
func (r *MySqlPostRepository) toOrderColumn(order string) (string, error) {
	if order == "" {
		return "id", nil
	}
	if order == "createdAt" {
		return "id", nil
	}
	if order == "location" {
		return "location_id", nil
	}
	if order == "userId" {
		return "user_id", nil
	}
	return "", ErrPostQueryInvalidOrderName
}

// モデルを変換
func (r *MySqlPostRepository) modelsToEntities(models []PostModel) ([]entity.Post, error) {
	imageMap := map[string]map[string]entity.Image{}
	tagMap := map[string]map[string]entity.Tag{}
	postMap := map[string]PostModel{}
	for _, model := range models {
		// Id変換
		id, err := uuid.FromBytes(model.Id)
		if err != nil {
			return nil, err
		}
		// タグマッピング
		if model.TagId != nil && model.TagName.Valid {
			tagId, err := uuid.FromBytes(*model.TagId)
			if err != nil {
				return nil, err
			}
			if _, ok := tagMap[id.String()]; !ok {
				tagMap[id.String()] = map[string]entity.Tag{}
			}
			if _, ok := tagMap[id.String()][tagId.String()]; !ok {
				tag, err := entity.NewTag(tagId, model.TagName.String)
				if err != nil {
					return nil, err
				}
				tagMap[id.String()][tagId.String()] = *tag
			}
		}
		// 画像マッピング
		if model.Image.Valid {
			if _, ok := imageMap[id.String()]; !ok {
				imageMap[id.String()] = map[string]entity.Image{}
			}
			if _, ok := imageMap[id.String()][model.Image.String]; !ok {
				image, err := entity.NewNameImage(model.Image.String)
				if err != nil {
					return nil, err
				}
				imageMap[id.String()][model.Image.String] = *image
			}
		}
		// 空投稿マッピング
		if _, ok := postMap[id.String()]; !ok {
			postMap[id.String()] = model
		}
	}
	// 投稿マッピング
	entities := []entity.Post{}
	for _, model := range models {
		// Id変換
		id, err := uuid.FromBytes(model.Id)
		if err != nil {
			return nil, err
		}
		// 追加済みか確認
		if _, ok := postMap[id.String()]; !ok {
			continue
		}
		// ユーザーId変換
		userId, err := uuid.FromBytes(model.UserId)
		if err != nil {
			return nil, err
		}
		// ロケーションId変換
		locationId, err := uuid.FromBytes(model.LocationId)
		if err != nil {
			return nil, err
		}
		// ロケーション作成
		location, err := entity.NewLocation(locationId, model.LocationName)
		if err != nil {
			return nil, err
		}
		// タグ変換
		tags := []entity.Tag{}
		if _, ok := tagMap[id.String()]; ok {
			for _, tag := range tagMap[id.String()] {
				tags = append(tags, tag)
			}
		}
		// 画像変換
		images := []entity.Image{}
		if _, ok := imageMap[id.String()]; ok {
			for _, image := range imageMap[id.String()] {
				images = append(images, image)
			}
		}
		// リストに投稿を追加
		post, err := entity.NewPost(id, userId, *location, model.Description, tags, images, model.CreatedAt, model.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entities = append(entities, *post)
		// マップから削除
		delete(postMap, id.String())
	}
	return entities, nil
}
