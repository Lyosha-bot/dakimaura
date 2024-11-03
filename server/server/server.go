package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"goserver/lib"
	"goserver/storage"
	"goserver/storage/database"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	port     string
	Database storage.Database
}

func NewServer(port string) (*Server, error) {

	credentialsDB := storage.Credentials{
		Host:     os.Getenv("DB_ENDPOINT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Catalog:  os.Getenv("DB_NAME"),
	}

	credentialsImages := storage.Credentials{
		Host:     os.Getenv("MINIO_ENDPOINT"),
		Username: os.Getenv("MINIO_USERNAME"),
		Password: os.Getenv("MINIO_PASSWORD"),
		Catalog:  os.Getenv("MINIO_BUCKET"),
	}

	dbClient, err := database.NewClient(credentialsDB, credentialsImages)
	if err != nil {
		return nil, lib.WrapErr("Minio init fail:", err)
	}

	return &Server{
		port:     port,
		Database: dbClient,
	}, nil
}

func (s *Server) Process() error {

	router := mux.NewRouter()

	router.HandleFunc("/add-product", s.addProduct).Methods("POST")
	router.HandleFunc("/get-product", s.getProduct).Methods("GET")
	router.HandleFunc("/get-category", s.getCategory).Methods("GET")

	log.Println("Server is up")

	err := http.ListenAndServe(":"+s.port, cors.AllowAll().Handler(router))
	return err
}

func (s *Server) addProduct(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("parse multipart form:", err)
		http.Error(w, "Couldn't parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("form file:", err)
		http.Error(w, "Couldn't get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	price, err := strconv.ParseUint(r.FormValue("price"), 10, 64)

	if err != nil {
		log.Println("price parse:", err)
		http.Error(w, "Couldn't parse uint", http.StatusBadRequest)
		return
	}

	data := storage.Product{
		Category:    r.FormValue("category"),
		Name:        r.FormValue("name"),
		Price:       price,
		Material:    r.FormValue("material"),
		Brand:       r.FormValue("brand"),
		ProduceTime: r.FormValue("produce_time"),
		Image:       header.Filename,
	}

	fileData := storage.FileData{
		File: file,
		Name: header.Filename,
		Size: header.Size,
	}

	err = s.Database.InsertProduct(data, fileData)
	if err != nil {
		log.Println("insert product:", err)
		http.Error(w, "Error inserting product:", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		log.Println("parse uint:", err)
		http.Error(w, "Parse uint:", http.StatusInternalServerError)
		return
	}

	product, err := s.Database.GetProduct(id)
	if err != nil {
		log.Println("get product:", err)
		http.Error(w, "Error encoding category:", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		log.Println("encode product:", err)
		http.Error(w, "Error encoding product:", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getCategory(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("category")
	products, err := s.Database.GetCategory(name)
	if err != nil {
		log.Println("get category:", err)
		http.Error(w, "Error getting category:", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Println("encode category:", err)
		http.Error(w, "Error encoding category:", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
