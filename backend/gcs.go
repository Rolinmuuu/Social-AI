package backend

import (
	"context"
	"fmt"
	"io"
	"socialai/constants"

	"cloud.google.com/go/storage"

)

var (
	GCSBackend GoogleCloudStorageBackendInterface
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

func InitGCSBackend() (GoogleCloudStorageBackendInterface, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	GCSBackend = &GoogleCloudStorageBackend{
		client: client,
		bucket: constants.GCS_BUCKET,
	}
	return GCSBackend, nil
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)
	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, r); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	// access control
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attribute, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	fmt.Printf("File is uploaded to GCS bucket %s\n", attribute.MediaLink)

	return attribute.MediaLink, err
}

func (backend *GoogleCloudStorageBackend) DeleteFromGCS(objectName string) error {
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)
	if err := object.Delete(ctx); err != nil {
		return err
	}
	return nil
}