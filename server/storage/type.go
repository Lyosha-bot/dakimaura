package storage

import (
	"io"
)

const (
	ProductInsertQuery      = "INSERT INTO products (category, name, price, material, brand, produce_time, image) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	ProductGetQuery         = "SELECT id, category, name, price, material, brand, produce_time, image FROM products WHERE id = $1 LIMIT 1"
	ProductGetCategoryQuery = "SELECT id, category, name, price, material, brand, produce_time, image FROM products WHERE category = $1"
	UserInsertQuery         = "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	UserGetQuery            = "SELECT id, username, password FROM users WHERE username = $1 LIMIT 1"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	ID          uint64 `json:"id"`
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
	Product(id uint64) (*Product, error)
	Category(category string) ([]Product, error)
	InsertUser(user User) error
	User(username string) (*User, error)
}
