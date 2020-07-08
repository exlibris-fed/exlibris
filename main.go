package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/exlibris-fed/exlibris/handler"
	"github.com/exlibris-fed/exlibris/handler/middleware"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatalf("HOST not provided")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT not provided")
	}
	if os.Getenv("DOMAIN") == "" {
		log.Fatalf("DOMAIN not provided")
	}
	if os.Getenv("SECRET") == "" {
		log.Fatalf("SECRET not provided")
	}

	db, err := gorm.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer db.Close()

	db.AutoMigrate(model.Author{})
	db.AutoMigrate(model.APObject{})
	db.AutoMigrate(model.Book{})
	db.AutoMigrate(model.OutboxEntry{})
	db.AutoMigrate(model.Read{})
	db.AutoMigrate(model.Review{})
	db.AutoMigrate(model.Subject{})
	db.AutoMigrate(model.User{})

	db.Model(&model.APObject{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&model.APObject{}).AddForeignKey("read_id", "reads(id)", "CASCADE", "CASCADE")

	db.Table("book_authors").AddForeignKey("author_id", "authors(id)", "CASCADE", "CASCADE")
	db.Table("book_authors").AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")

	db.Table("book_subjects").AddForeignKey("subject_id", "subjects(id)", "CASCADE", "CASCADE")
	db.Table("book_subjects").AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")

	db.Model(&model.OutboxEntry{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.Model(&model.Read{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")
	db.Model(&model.Read{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.Model(&model.Review{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&model.Review{}).AddForeignKey("book_id", "books(id)", "CASCADE", "CASCADE")

	h := handler.New(db)
	m := middleware.New(db)

	r := mux.NewRouter()
	r.Use(m.ExtractUsername)

	// APIs
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", h.Register).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/authenticate", h.Authenticate).Methods(http.MethodPost, http.MethodOptions)

	books := api.PathPrefix("/book").Subrouter()
	books.Use(m.Authenticated)
	books.HandleFunc("", h.SearchBooks)
	books.HandleFunc("/{book}/read", h.Read).Methods(http.MethodPost, http.MethodOptions)
	books.HandleFunc("/read", h.GetReads)

	// inbox/outbox handle authentication as part of the go-fed flow. ExtractUsername will populate it if present.
	api.Handle("/user/{username}/inbox", m.WithUserModel(http.HandlerFunc(h.HandleInbox)))
	api.Handle("/user/{username}/outbox", m.WithUserModel(http.HandlerFunc(h.HandleOutbox)))
	api.Handle("/@{username}/inbox", m.WithUserModel(http.HandlerFunc(h.HandleInbox)))
	api.Handle("/@{username}/outbox", m.WithUserModel(http.HandlerFunc(h.HandleOutbox)))

	// App
	r.HandleFunc("/.well-known/acme-challenge/{id}", h.HandleChallenge)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))
	corsRouter := handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Access-Control-Allow-Origin"}))
	loggedRouter := handlers.LoggingHandler(os.Stdout, corsRouter(r))

	addr := net.JoinHostPort(host, port)
	log.Println("Starting on", addr)

	server := &http.Server{
		Handler:      loggedRouter,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
