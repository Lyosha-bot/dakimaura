package storage

import "io"

const (
	ProductInsertQuery = "INSERT INTO products (category, name, price, material, brand, produce_time, image) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	ProductGetQuery    = "SELECT * FROM products WHERE id = $1 LIMIT 1"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	Category    string `json:"category"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Material    string `json:"material"`
	Brand       string `json:"brand"`
	ProduceTime string `json:"produce_time"`
	Image       string `json:"image"`
}

type Images interface {
	UploadImage(file io.Reader, name string, size int64) error
	LoadImage(name string) error
	DeleteImage(name string) error
}

type Database interface {
	InsertProduct(product Product) error
	GetProduct(id uint) (Product, error)
	GetCategory(category string) ([]Product, error)
}
