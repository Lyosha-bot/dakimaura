package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"goserver/lib"
	"goserver/storage"
	"goserver/storage/images"
	"log"
	"time"
)

type Client struct {
	pool   *pgxpool.Pool
	Images storage.Images
}

func createContext() (context.Context, context.CancelFunc) { // Бесполезно, но удобно
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func NewClient(credentialsDB, credentialsImages storage.Credentials) (storage.Database, error) {
	ctx, cancel := createContext()
	defer cancel()

	//Images

	imagesClient, err := images.NewClient(credentialsImages)
	if err != nil {
		return nil, lib.WrapErr("minio init:", err)
	}

	//Database

	url := fmt.Sprintf("postgres://%s:%s@%s/%s", credentialsDB.Username, credentialsDB.Password, credentialsDB.Host, credentialsDB.Catalog)

	log.Println(url)

	pool, err := pgxpool.New(ctx, url)

	if err != nil {
		return nil, lib.WrapErr("database pool:", err)
	}

	return &Client{pool: pool, Images: imagesClient}, nil
}

func (c *Client) InsertProduct(product storage.Product, imageData storage.FileData) error {
	ctx, cancel := createContext()
	defer cancel()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return lib.WrapErr("acquire:", err)
	}
	defer conn.Release()

	err = c.Images.UploadImage(imageData)
	if err != nil {
		return lib.WrapErr("insert product:", err)
	}

	product.Image = imageData.Name

	row := conn.QueryRow(ctx,
		storage.ProductInsertQuery,
		product.Category, product.Name, product.Price, product.Material, product.Brand, product.ProduceTime, product.Image)

	var id uint
	err = row.Scan(&id)

	if err != nil {
		_ = c.Images.DeleteImage(imageData.Name)
		return lib.WrapErr("insert:", err)
	}

	return nil
}

func (c *Client) GetProduct(id uint64) (*storage.Product, error) {
	ctx, cancel := createContext()
	defer cancel()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.WrapErr("acquire:", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, storage.ProductGetQuery, id)

	var data storage.Product
	err = row.Scan(&data.ID, &data.Category, &data.Name, &data.Price, &data.Material, &data.Brand, &data.ProduceTime, &data.Image)

	if err != nil {
		return nil, lib.WrapErr("row scan:", err)
	}

	return &data, lib.WrapIfErr("load image:", err)
}

func (c *Client) GetCategory(category string) ([]storage.Product, error) {
	ctx, cancel := createContext()
	defer cancel()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.WrapErr("acquire:", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, storage.ProductGetCategoryQuery, category)
	if err != nil {
		return nil, lib.WrapErr("category query:", err)
	}

	res := make([]storage.Product, 0, 1)

	var data storage.Product
	for rows.Next() {
		err = rows.Scan(&data.ID, &data.Category, &data.Name, &data.Price, &data.Material, &data.Brand, &data.ProduceTime, &data.Image)
		if err != nil {
			return nil, lib.WrapErr("row scan:", err)
		}

		res = append(res, data)
	}
	return res, nil
}
