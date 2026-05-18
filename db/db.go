package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

func GetActiveDatabaseSession() *gorm.DB {
	return db
}

func GetDatabaseConnection(username string, password string, dbName string) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	// Connect to default 'postgres' database as your existing user
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=postgres port=5432 sslmode=disable", username, password)

	adminDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	defer adminDB.Close()

	// Check if database exists
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	adminDB.QueryRow(query).Scan(&exists)

	if !exists {
		// Create database
		createSQL := fmt.Sprintf("CREATE DATABASE %s", dbName)
		_, err = adminDB.Exec(createSQL)
		if err != nil {
			return nil, fmt.Errorf("[-] failed to create database: %w", err)
		}
		log.Printf("[+] Database '%s' created", dbName)
	} else {
		log.Printf("[+] Database '%s' already exists", dbName)
	}

	// Connect to your database with GORM
	appDSN := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable", username, password, dbName)
	db, err = gorm.Open(postgres.Open(appDSN), &gorm.Config{SkipDefaultTransaction: true})
	return db, err
}
