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
	"github.com/exlibris-fed/exlibris/infrastructure"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	cfg := config.Load()

	db := infrastructure.New(cfg.DSN)
	defer db.Close()

	infrastructure.Migrate(db)

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
	ap.PathPrefix("/").Handler(http.HandlerFunc(h.HandleActivityPubAction))

	r.PathPrefix("/user/").Handler(http.HandlerFunc(h.HandleActivityPubAction))

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
