package server

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"goserver/lib"
	"goserver/storage"
	"goserver/storage/database"
	"goserver/storage/images"
	"log"
	"net/http"
	"os"
)

type Server struct {
	port     string
	Images   storage.Images
	Database storage.Database
}

func NewServer(port string) (*Server, error) {
	imagesClient, err := images.NewClient(os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_USER"), os.Getenv("MINIO_PASSWORD"))
	if err != nil {
		return nil, lib.WrapErr("Minio init fail:", err)
	}

	dbClient, err := database.NewClient(os.Getenv("DB_ENDPOINT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	if err != nil {
		return nil, lib.WrapErr("Minio init fail:", err)
	}

	return &Server{
		port:     port,
		Images:   imagesClient,
		Database: dbClient,
	}, nil
}

func (s *Server) Process() error {

	router := mux.NewRouter()

	router.HandleFunc("/add-product", s.addProduct).Methods("POST")

	log.Println("Server is up")

	//c := cors.New(cors.Options{
	//	AllowedOrigins: []string{"localhost:8073", "nginx:8073"},
	//	AllowedMethods: []string{"GET", "POST"},
	//})

	err := http.ListenAndServe(":"+s.port, cors.Default().Handler(router))
	return err
}

func (s *Server) addProduct(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Couldn't parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Couldn't get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data := storage.Product{
		Category:    r.FormValue("category"),
		Name:        r.FormValue("name"),
		Material:    r.FormValue("material"),
		Brand:       r.FormValue("brand"),
		ProduceTime: r.FormValue("produce_time"),
		Image:       header.Filename,
	}

	err = s.Images.UploadImage(file, header.Filename, header.Size)
	if err != nil {
		log.Println("error uploading image:", err)
		http.Error(w, "Error uploading image:", http.StatusInternalServerError)
		return
	}

	err = s.Database.InsertProduct(data)
	if err != nil {
		_ = s.Images.DeleteImage(header.Filename)
		log.Println("error inserting product:", err)
		http.Error(w, "Error inserting product:", http.StatusInternalServerError)
		return
	}
}
