package server

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"goserver/lib"
	"goserver/storage"
	"goserver/storage/database"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const tokenDuration = time.Hour * 24

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
		return nil, lib.WrapErr("database init fail:", err)
	}

	return &Server{
		port:     port,
		Database: dbClient,
	}, nil
}

func (s *Server) Process() error {

	router := mux.NewRouter()

	router.HandleFunc("/add-product", s.addProduct).Methods("POST")
	router.HandleFunc("/get-product", s.product).Methods("GET")
	router.HandleFunc("/get-category", s.category).Methods("GET")
	router.HandleFunc("/register", s.register).Methods("POST")
	router.HandleFunc("/login", s.login).Methods("POST")
	router.HandleFunc("/fetch", s.fetch).Methods("GET")
	router.HandleFunc("/logout", s.logout).Methods("POST")

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

func (s *Server) product(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		log.Println("parse uint:", err)
		http.Error(w, "Parse uint:", http.StatusInternalServerError)
		return
	}

	product, err := s.Database.Product(id)
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

func (s *Server) category(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("category")
	products, err := s.Database.Category(name)
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

func invalidUser(user storage.User) error {
	if user.Username == "" {
		return errors.New("no username")
	}

	if len(user.Username) < 5 {
		return errors.New("too short username")
	}

	if user.Password == "" {
		return errors.New("no password")
	}

	if len(user.Password) < 5 {
		return errors.New("too short password")
	}

	return nil
}

func newToken(user *storage.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", lib.WrapErr("token sign:", err)
	}
	return tokenString, nil
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("parse multipart form:", err)
		http.Error(w, "Форма неверно заполнена", http.StatusBadRequest)
		return
	}

	data := storage.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if err := invalidUser(data); err != nil {
		log.Println("invalid user data:", err)
		http.Error(w, "Данные неверно заполнены", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("hash password:", err)
		http.Error(w, "Ошибка сервера, попробуйте еще раз", http.StatusInternalServerError)
		return
	}

	data.Password = string(hash)

	err = s.Database.InsertUser(data)
	if err != nil {
		log.Println("insert user:", err)
		http.Error(w, "Ошибка сервера, попробуйте еще раз", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("parse multipart form:", err)
		http.Error(w, "Форма неверно заполнена", http.StatusBadRequest)
		return
	}

	formData := storage.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if err := invalidUser(formData); err != nil {
		log.Println("invalid user data:", err)
		http.Error(w, "Данные неверно заполнены", http.StatusBadRequest)
		return
	}

	user, err := s.Database.User(formData.Username)
	if err != nil {
		log.Println("incorrect user:", err)
		http.Error(w, "Неверный логин или пароль", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formData.Password))
	if err != nil {
		log.Println("incorrect password:", err)
		http.Error(w, "Неверный логин или пароль", http.StatusBadRequest)
		return
	}

	token, err := newToken(user, os.Getenv("APP_SECRET"), tokenDuration)
	if err != nil {
		log.Println("token generation:", err)
		http.Error(w, "Ошибка сервера, попробуйте еще раз", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(tokenDuration),
		Secure:   false,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth_token")
	if err != nil {
		log.Println("get cookie:", err)
		http.Error(w, "Нет токена", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   false,
		HttpOnly: true,
	})

	log.Println("DELETED COOKIE")

	w.WriteHeader(http.StatusOK)
}

func (s *Server) fetch(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		log.Println("get cookie:", err)
		http.Error(w, "Нет токена", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("APP_SECRET")), nil
	})

	if err != nil {
		log.Println("parse token:", err)
		http.Error(w, "Неверный токен", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("token claims:", err)
		http.Error(w, "Неверный токен", http.StatusBadRequest)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		log.Println("token username:", err)
		http.Error(w, "Неверный токен", http.StatusBadRequest)
		return
	}

	_, ok = claims["id"]
	if !ok {
		log.Println("token user id:", err)
		http.Error(w, "Неверный токен", http.StatusBadRequest)
		return
	}

	w.Write([]byte(username))
	w.WriteHeader(http.StatusOK)
}
