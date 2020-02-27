// Package model contains the models used by exlibris. Each model lives in its own file and should register itself via the registerModel function so that migrations will be applied.
package model

import (
	"github.com/jinzhu/gorm"
)

var modelList []interface{}

func init() {
	modelList = []interface{}{}
}

func registerModel(m interface{}) {
	modelList = append(modelList, m)
}

// ApplyMigrations is used to migrate all models. If you only wish to migrate one model, call db.AutoMigrate on it individually.
func ApplyMigrations(db *gorm.DB) error {
	for _, m := range modelList {
		db.AutoMigrate(m)
	}

	// TODO does AutoMigrate return an error?
	return nil
}
