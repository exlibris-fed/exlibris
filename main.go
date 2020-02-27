package main

import (
    "fmt"
    "time"
    "log"
    "net/http"
    "os"

    "github.com/exlibris-fed/exlibris/activitypub"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/gorilla/mux"
)

func main() {
    log.Println(os.Getenv("POSTGRES_CONNECTION"))
    db, err := gorm.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
    if err != nil {
        log.Fatalf("unable to connect to database: %s", err)
    }
    defer db.Close()

    ap := activitypub.New(db)

    r := mux.NewRouter()
    r.HandleFunc("/{user}/inbox", ap.HandleInbox)
    r.HandleFunc("/{user}/outbox", ap.HandleOutbox)

    server := &http.Server{
        Handler: r,
        Addr: "127.0.0.1:8080", // TODO config it somewhere?
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    fmt.Println("exlibris running")
    log.Fatal(server.ListenAndServe())
}
