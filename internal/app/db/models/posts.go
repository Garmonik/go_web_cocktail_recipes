package models

import "time"

type Image struct {
	ID   uint   `gorm:"primaryKey"`
	Path string `gorm:"not null;unique"`
}

type Post struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Image       Image     `gorm:"foreignKey:ImageID"`
	Author      User      `gorm:"foreignKey:AuthorID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	AuthorID    uint
	ImageID     uint
}

type Like struct {
	ID       uint `gorm:"primaryKey"`
	Author   User `gorm:"foreignKey:AuthorID"`
	Post     Post `gorm:"foreignKey:PostID"`
	PostID   uint
	AuthorID uint
}

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	text      string    `gorm:"not null"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	AuthorID  uint
}
