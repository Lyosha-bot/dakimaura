package images

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"goserver/lib"
	"goserver/storage"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"time"
)

const bucketName = "images"

type Client struct {
	client *minio.Client
}

func createContext() (context.Context, context.CancelFunc) { // Бесполезно, но удобно
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func NewClient(endpoint, username, password string) (storage.Images, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, lib.WrapErr("minio connection failed: ", err)
	}

	ctx, cancel := createContext()
	defer cancel()

	found, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, lib.WrapErr("minio bucket check failed: ", err)
	}

	if !found {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, lib.WrapErr("minio bucket check failed: ", err)
		}
	}

	return &Client{client: client}, nil
}

func (c *Client) UploadImage(file io.Reader, name string, size int64) error {
	ext := filepath.Ext(name)
	contentType := mime.TypeByExtension(ext)

	ctx, cancel := createContext()
	defer cancel()

	_, err := c.client.PutObject(ctx, bucketName, name, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (c *Client) LoadImage(name string) error {
	ctx, cancel := createContext()
	defer cancel()

	_, err := c.client.PresignedGetObject(ctx, bucketName, name, 60*60*24, make(url.Values))
	return err
}

func (c *Client) DeleteImage(name string) error {
	ctx, cancel := createContext()
	defer cancel()

	err := c.client.RemoveObject(ctx, bucketName, name, minio.RemoveObjectOptions{
		ForceDelete:      true,
		GovernanceBypass: true,
	})
	return err
}
