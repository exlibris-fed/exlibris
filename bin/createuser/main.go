// DO NOT USE THIS IN A PRODUCTION ENVIRONMENT! JEEZY CREEZY!
package main

import (
	"fmt"
	"os"

	"github.com/exlibris-fed/exlibris/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: go run . username displayname password")
		os.Exit(1)
	}

	u := model.User{
		Username:    os.Args[1],
		DisplayName: os.Args[2],
	}
	u.SetPassword(os.Args[3])
	err := u.GenerateKeys()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := gorm.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
	if err != nil {
		fmt.Printf("unable to connect to database: %s", err)
		os.Exit(1)
	}
	defer db.Close()

	db.Create(&u)
}
