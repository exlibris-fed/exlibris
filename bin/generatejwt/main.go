// DO NOT USE THIS IN A PRODUCTION ENVIRONMENT! JEEZY CREEZY!
package main

import (
	"fmt"
	"os"

	"github.com/exlibris-fed/exlibris/key"
	"github.com/exlibris-fed/exlibris/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . username")
		os.Exit(1)
	}

	var u model.User

	db, err := gorm.Open("postgres", os.Getenv("POSTGRES_CONNECTION"))
	if err != nil {
		fmt.Printf("unable to connect to database: %s", err)
		os.Exit(1)
	}
	defer db.Close()

	db.First(&u, "username = ?", os.Args[1])
	if len(u.Password) == 0 {
		fmt.Printf("user '%s' does not exist", os.Args[1])
		os.Exit(1)
	}

	k, err := key.DeserializeRSAPrivateKey(u.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	jwt, err := key.GenerateJWT(k)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jwt)
}
