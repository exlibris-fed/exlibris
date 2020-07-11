package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/exlibris-fed/exlibris/config"
	"github.com/exlibris-fed/exlibris/handler"
	"github.com/exlibris-fed/exlibris/handler/middleware"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open("postgres", cfg.DSN)
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
	db.AutoMigrate(model.RegistrationKey{})
	db.AutoMigrate(model.Cover{})

	db.Model(&model.APObject{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&model.APObject{}).AddForeignKey("read_id", "reads(id)", "CASCADE", "CASCADE")

	db.Table("book_authors").AddForeignKey("author_open_library_id", "authors(open_library_id)", "CASCADE", "CASCADE")
	db.Table("book_authors").AddForeignKey("book_open_library_id", "books(open_library_id)", "CASCADE", "CASCADE")

	db.Table("book_subjects").AddForeignKey("subject_id", "subjects(id)", "CASCADE", "CASCADE")
	db.Table("book_subjects").AddForeignKey("book_open_library_id", "books(open_library_id)", "CASCADE", "CASCADE")

	db.Table("book_covers").AddForeignKey("book_open_library_id", "books(open_library_id)", "CASCADE", "CASCADE")
	db.Table("book_covers").AddForeignKey("cover_id", "covers(id)", "CASCADE", "CASCADE")

	db.Model(&model.OutboxEntry{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.Model(&model.Read{}).AddForeignKey("book_id", "books(open_library_id)", "CASCADE", "CASCADE")
	db.Model(&model.Read{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.Model(&model.Review{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&model.Review{}).AddForeignKey("book_id", "books(open_library_id)", "CASCADE", "CASCADE")
	db.Model(&model.RegistrationKey{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	h := handler.New(db, cfg)
	m := middleware.New(db)

	r := mux.NewRouter()
	r.Use(m.ExtractUsername)

	// APIs
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", h.Register).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/authenticate", h.Authenticate).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/verify/resend/{user}", h.ResendVerificationKey).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/verify/{key}", h.VerifyKey).Methods(http.MethodGet, http.MethodOptions)
	api.Handle("/user/{username}", http.HandlerFunc(h.HandleActivityPubProfile))

	books := api.PathPrefix("/book").Subrouter()
	books.Use(m.WithUserModel)
	books.HandleFunc("", h.SearchBooks).Methods(http.MethodGet, http.MethodOptions)
	books.HandleFunc("/{book}/read", h.Read).Methods(http.MethodPost, http.MethodOptions)
	books.HandleFunc("/read", h.GetReads).Methods(http.MethodGet, http.MethodOptions)
	books.HandleFunc("/{book}", h.GetBook).Methods(http.MethodGet, http.MethodOptions)
	books.HandleFunc("/{book}/review", h.Review).Methods(http.MethodPost, http.MethodOptions, http.MethodGet)

	// inbox/outbox handle authentication as part of the go-fed flow. ExtractUsername will populate it if present.
	ap := r.Headers("Accept", "application/activity+json").Subrouter()
	// TODO add withusermodel here?
	ap.Handle("/user/{username}", http.HandlerFunc(h.HandleActivityPubProfile))
	ap.Handle("/user/{username}/inbox", m.WithUserModel(http.HandlerFunc(h.HandleInbox)))
	ap.Handle("/user/{username}/outbox", m.WithUserModel(http.HandlerFunc(h.HandleOutbox)))
	ap.Handle("/@{username}", http.HandlerFunc(h.HandleActivityPubProfile))
	ap.Handle("/@{username}/inbox", m.WithUserModel(http.HandlerFunc(h.HandleInbox)))
	ap.Handle("/@{username}/outbox", m.WithUserModel(http.HandlerFunc(h.HandleOutbox)))

	// JSON handlers. may not be needed? Hackathon!!
	jsonRouter := r.Headers("Accept", "application/json").Subrouter()
	// TODO add withusermodel here?
	jsonRouter.Handle("/user/{username}", http.HandlerFunc(h.HandleActivityPubProfile))
	jsonRouter.Handle("/@{username}", http.HandlerFunc(h.HandleActivityPubProfile))

	// App
	r.HandleFunc("/.well-known/acme-challenge/{id}", h.HandleChallenge)
	r.HandleFunc("/.well-known/webfinger", h.HandleWebfinger)
	r.PathPrefix("/").Handler(http.HandlerFunc(h.HandleStaticFile))
	corsRouter := handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Access-Control-Allow-Origin"}))
	loggedRouter := handlers.LoggingHandler(os.Stdout, corsRouter(r))

	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	log.Println("Starting on", addr)

	server := &http.Server{
		Handler:      loggedRouter,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
