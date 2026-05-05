package test

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Users []User `gorm:"many2many:user_languages;"`
}

// Custom join table
type UserLanguage struct {
	UserID      uint   `gorm:"primaryKey"`
	LanguageID  uint   `gorm:"primaryKey"`
	Proficiency string `gorm:"size:20"` // Extra field
	CreatedAt   time.Time
}
