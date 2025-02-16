package models

type Image struct {
	ID   uint   `gorm:"primaryKey"`
	Path string `gorm:"not null;unique"`
}

type Post struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`

	ImageID uint
	Image   Avatar `gorm:"foreignKey:ImageID"`

	AuthorID uint
	Author   User `gorm:"foreignKey:AuthorID"`
}

type Like struct {
	ID uint `gorm:"primaryKey"`

	AuthorID uint
	Author   User `gorm:"foreignKey:AuthorID"`

	PostID uint
	Post   Post `gorm:"foreignKey:PostID"`
}

type Comment struct {
	ID   uint   `gorm:"primaryKey"`
	text string `gorm:"not null"`

	AuthorID uint
	Author   User `gorm:"foreignKey:AuthorID"`
}
