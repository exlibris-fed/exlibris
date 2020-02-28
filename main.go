package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/exlibris-fed/exlibris/handler"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	db, err := gorm.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer db.Close()

	db.AutoMigrate(model.Author{})
	db.AutoMigrate(model.Book{})
	db.AutoMigrate(model.BookAuthor{})
	db.AutoMigrate(model.BookSubject{})
	db.AutoMigrate(model.Read{})
	db.AutoMigrate(model.Review{})
	db.AutoMigrate(model.Subject{})
	db.AutoMigrate(model.User{})

	h := handler.New(db)

	r := mux.NewRouter()
	r.HandleFunc("/register", h.Register).Methods(http.MethodPost)
	r.HandleFunc("/authenticate", h.Authenticate).Methods(http.MethodPost)
	r.HandleFunc("/book", h.SearchBooks)
	r.HandleFunc("/book/{book}/read", h.Read).Methods(http.MethodPost)
	r.HandleFunc("/user/{username}/inbox", h.HandleInbox)
	r.HandleFunc("/user/{username}/outbox", h.HandleOutbox)
	r.HandleFunc("/@{username}/inbox", h.HandleInbox)
	r.HandleFunc("/@{username}/outbox", h.HandleOutbox)
	r.HandleFunc("/fedtest", h.FederationTest).Methods(http.MethodPost)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	addr := net.JoinHostPort(host, port)
	log.Println("Starting on ", addr)

	server := &http.Server{
		Handler:      handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(loggedRouter),
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("exlibris running")
	log.Fatal(server.ListenAndServe())
}
