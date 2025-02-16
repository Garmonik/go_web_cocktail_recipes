package db

import (
	"fmt"
	models2 "github.com/Garmonik/go_web_cocktail_recipes/internal/app/db/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DataBase struct {
	db *gorm.DB
}

func New(storagePath string) (*DataBase, error) {
	const operation = "qs.sqlite.New"

	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	database := &DataBase{db: db}

	if err := database.migrate(); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	return database, nil
}

// migrate
func (d *DataBase) migrate() error {
	return d.db.AutoMigrate(
		&models2.Avatar{},
		&models2.User{},
		&models2.Post{},
		&models2.Image{},
		&models2.Comment{},
		&models2.Like{},
	)
}
