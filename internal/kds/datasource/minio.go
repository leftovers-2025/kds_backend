package datasource

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetMinIOClient() *minio.Client {
	client, err := minio.New(
		getMinioEndpoint(),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				getMinioAccessKeyId(),
				getMinioSecret(),
				"",
			),
			Secure: isMinioUseSSL(),
		},
	)
	if err != nil {
		panic("failed to get minio client.\n" + err.Error())
	}
	err = MigrateMinIO(client)
	if err != nil {
		panic("MinIO Migration Error: " + err.Error())
	}
	return client
}

func MigrateMinIO(client *minio.Client) error {
	imageBucketName := "images"
	exists, err := client.BucketExists(context.Background(), imageBucketName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	err = client.MakeBucket(context.Background(), imageBucketName, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	return nil
}

func getMinioEndpoint() string {
	endpoint, ok := os.LookupEnv("MINIO_ENDPOINT")
	if !ok {
		panic("\"MINIO_ENDPOINT\" is not set")
	}
	return endpoint
}

func isMinioUseSSL() bool {
	endpoint, ok := os.LookupEnv("MINIO_USE_SSL")
	if !ok {
		panic("\"MINIO_ENDPOINT\" is not set")
	}
	if endpoint == "1" {
		return true
	}
	return false
}

func getMinioAccessKeyId() string {
	accessKey, ok := os.LookupEnv("MINIO_ROOT_USER")
	if !ok {
		panic("\"MINIO_ROOT_USER\" is not set")
	}
	return accessKey
}

func getMinioSecret() string {
	accessKey, ok := os.LookupEnv("MINIO_ROOT_PASSWORD")
	if !ok {
		panic("\"MINIO_ROOT_PASSWORD\" is not set")
	}
	return accessKey
}
