package test

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Product struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"size:200;not null;index"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
	Stock     int     `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	// Connection string
	dsn := "host=localhost user=postgres password=yourpass dbname=testdb port=5432 sslmode=disable TimeZone=UTC"

	// Connect with logger
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate schema
	db.AutoMigrate(&Product{})

	// Create
	product := Product{Name: "Laptop", Price: 999.99, Stock: 10}
	db.Create(&product)

	// Read
	var p Product
	db.First(&p, product.ID)
	fmt.Printf("Found product: %+v\n", p)

	// Update
	db.Model(&p).Update("Price", 899.99)

	// Query all
	var products []Product
	db.Find(&products)
	fmt.Printf("All products: %+v\n", products)

	// Delete
	db.Delete(&p)
}
