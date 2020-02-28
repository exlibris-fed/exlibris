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
	host := os.Getenv("APP_HOST")
	if host == "" {
		log.Fatalf("APP_HOST not provided")
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatalf("APP_PORT not provided")
	}
	db, err := gorm.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer db.Close()

	model.ApplyMigrations(db)
	h := handler.New(db)

	r := mux.NewRouter()
	r.HandleFunc("/book", h.SearchBooks)
	r.HandleFunc("/user/{username}/inbox", h.HandleInbox)
	r.HandleFunc("/user/{username}/outbox", h.HandleOutbox)
	r.HandleFunc("/@{username}/inbox", h.HandleInbox)
	r.HandleFunc("/@{username}/outbox", h.HandleOutbox)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	addr := net.JoinHostPort(host, port)
	log.Println("Starting on ", addr)

	server := &http.Server{
		Handler:      loggedRouter,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("exlibris running")
	log.Fatal(server.ListenAndServe())
}
