package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/exlibris-fed/exlibris/activitypub"
	"github.com/exlibris-fed/exlibris/handler"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	fmt.Println(os.Getenv("POSTGRES_CONNECTION"))
	conn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
	if err != nil {
		log.Fatal(err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatalf("Could not ping: %s", err)
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

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
