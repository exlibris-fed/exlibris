package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/exlibris-fed/exlibris/activitypub"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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

	r := mux.NewRouter()
	r.HandleFunc("/{username}/inbox", ap.HandleInbox)
	r.HandleFunc("/{username}/outbox", ap.HandleOutbox)

	server := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080", // TODO config it somewhere?
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("exlibris running")
	log.Fatal(server.ListenAndServe())
}
