package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/exlibris-fed/exlibris/activitypub"
	"github.com/exlibris-fed/exlibris/handler"
	"github.com/exlibris-fed/exlibris/model"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	db, err := gorm.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}
	defer db.Close()

	model.ApplyMigrations(db)
	ap := activitypub.New(db)
	h := handler.New(db)

	r := mux.NewRouter()
	r.HandleFunc("/book", h.SearchBooks)
	r.HandleFunc("/user/{username}/inbox", ap.HandleInbox)
	r.HandleFunc("/user/{username}/outbox", ap.HandleOutbox)
	r.HandleFunc("/@{username}/inbox", ap.HandleInbox)
	r.HandleFunc("/@{username}/outbox", ap.HandleOutbox)

	addr := net.JoinHostPort(os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))

	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("exlibris running")
	log.Fatal(server.ListenAndServe())
}
