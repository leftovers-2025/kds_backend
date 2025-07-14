package s3

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

const (
	S3_IMAGES_BUCKET = "images"
)

type S3Repository struct {
	client *minio.Client
}

func NewS3Repository(client *minio.Client) *S3Repository {
	if client == nil {
		panic("nil MinIO Client")
	}
	return &S3Repository{
		client: client,
	}
}

// 画像をアップロード
func (r *S3Repository) UploadImage(name string, fileHeader *multipart.FileHeader) error {
	// ファイルを開く
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// アップロード
	_, err = r.client.PutObject(context.Background(), S3_IMAGES_BUCKET, name, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	return err
}
