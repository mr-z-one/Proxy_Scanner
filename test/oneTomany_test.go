package test

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Posts []Post `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Post struct {
	ID      uint `gorm:"primaryKey"`
	Title   string
	Content string
	UserID  uint `gorm:"index"` // Foreign key
}
