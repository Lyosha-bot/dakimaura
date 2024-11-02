package storage

import "io"

const (
	ProductInsertQuery      = "INSERT INTO products (category, name, price, material, brand, produce_time, image) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	ProductGetQuery         = "SELECT * FROM products WHERE id = $1 LIMIT 1"
	ProductGetCategoryQuery = "SELECT * FROM products WHERE category = $1"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	Category    string `json:"category"`
	Name        string `json:"name"`
	Price       uint64 `json:"price"`
	Material    string `json:"material"`
	Brand       string `json:"brand"`
	ProduceTime string `json:"produce_time"`
	Image       string `json:"image"`
}

type Credentials struct {
	Host     string
	Username string
	Password string
	Catalog  string
}

type FileData struct {
	File io.Reader
	Name string
	Size int64
}

type Images interface {
	UploadImage(imageData FileData) error
	DeleteImage(name string) error
}

type Database interface {
	InsertProduct(product Product, imageData FileData) error
	GetProduct(id uint64) (*Product, error)
	GetCategory(category string) ([]Product, error)
}
