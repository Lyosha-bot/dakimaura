package images

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"goserver/lib"
	"goserver/storage"
	"mime"
	"path/filepath"
	"time"
)

type Client struct {
	client     *minio.Client
	bucketName string
}

func createContext() (context.Context, context.CancelFunc) { // Бесполезно, но удобно
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func NewClient(credentialsImages storage.Credentials) (storage.Images, error) {
	client, err := minio.New(credentialsImages.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(credentialsImages.Username, credentialsImages.Password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, lib.WrapErr("minio connection failed: ", err)
	}

	ctx, cancel := createContext()
	defer cancel()

	found, err := client.BucketExists(ctx, credentialsImages.Catalog)
	if err != nil {
		return nil, lib.WrapErr("minio bucket check failed: ", err)
	}

	if !found {
		err = client.MakeBucket(ctx, credentialsImages.Catalog, minio.MakeBucketOptions{})
		if err != nil {
			return nil, lib.WrapErr("minio bucket check failed: ", err)
		}
	}

	return &Client{client: client, bucketName: credentialsImages.Catalog}, nil
}

func (c *Client) UploadImage(imageData storage.FileData) error {
	ext := filepath.Ext(imageData.Name)
	contentType := mime.TypeByExtension(ext)

	ctx, cancel := createContext()
	defer cancel()

	_, err := c.client.PutObject(ctx, c.bucketName, imageData.Name, imageData.File, imageData.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) DeleteImage(name string) error {
	ctx, cancel := createContext()
	defer cancel()

	err := c.client.RemoveObject(ctx, c.bucketName, name, minio.RemoveObjectOptions{
		ForceDelete:      true,
		GovernanceBypass: true,
	})
	return err
}
