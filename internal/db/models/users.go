package models

type Avatar struct {
	ID   uint   `gorm:"primaryKey"`
	Path string `gorm:"not null;unique"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null;uniqueIndex"`
	Password string `gorm:"not null"`
	Bio      string `gorm:"type:text"`
	AvatarID uint
	Avatar   Avatar `gorm:"foreignKey:AvatarID"`
}
