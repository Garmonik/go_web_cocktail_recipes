package db

import (
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/db/models"
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
		&models.Avatar{},
		&models.User{},
		&models.Post{},
		&models.Image{},
		&models.Comment{},
		&models.Like{},
	)
}
