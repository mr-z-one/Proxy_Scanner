package test

type User struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Profile Profile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Profile struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex"` // Foreign key
	Bio    string
	Avatar string
}
