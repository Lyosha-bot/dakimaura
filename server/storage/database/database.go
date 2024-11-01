package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"goserver/lib"
	"goserver/storage"
	"log"
	"time"
)

type Client struct {
	pool *pgxpool.Pool
}

func createContext() (context.Context, context.CancelFunc) { // Бесполезно, но удобно
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func NewClient(endpoint, username, password, database string) (*Client, error) {
	ctx, cancel := createContext()
	defer cancel()

	url := fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, endpoint, database)

	log.Println(url)

	pool, err := pgxpool.New(ctx, url)

	if err != nil {
		return nil, lib.WrapErr("database pool error:", err)
	}

	return &Client{pool: pool}, nil
}

func (c *Client) InsertProduct(product storage.Product) error {
	ctx, cancel := createContext()
	defer cancel()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return lib.WrapErr("acquire error:", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx,
		storage.ProductInsertQuery,
		product.Category, product.Name, product.Price, product.Material, product.Brand, product.ProduceTime, product.Image)

	var id uint
	err = row.Scan(&id)
	return lib.WrapIfErr("insert error:", err)
}

func (c *Client) GetProduct(id uint) (*storage.Product, error) {
	ctx, cancel := createContext()
	defer cancel()

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, lib.WrapErr("acquire error:", err)
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, storage.ProductGetQuery, id)

	var data storage.Product
	err = row.Scan(&id, &data.Category, &data.Name, &data.Price, &data.Material, &data.Brand, &data.ProduceTime, &data.Image)

}

func (c *Client) GetCategory(category string) ([]storage.Product, error) {

}
