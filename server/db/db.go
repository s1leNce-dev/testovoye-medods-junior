package db

import (
	"fmt"
	"log"
	"medods/models"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

func InitDatabase() {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			os.Getenv("PGSQL_HOST"), os.Getenv("PGSQL_PORT"), os.Getenv("PGSQL_USER"), os.Getenv("PGSQL_PASSWORD"), os.Getenv("PGSQL_DB"))

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("[FATAL] [PSQL] %s", err.Error())
		}

		if err := db.AutoMigrate(&models.User{}, &models.RefreshTokenSessions{}); err != nil {
			log.Fatalf("[FATAL] [PSQL] %s", err.Error())
		}

		log.Println("[INFO] [PSQL] Connected to PostgreSQL")
	})
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	if sqlDB, err := db.DB(); err != nil {
		log.Fatalf("[FATAL] [PSQL] %s", err.Error())
	} else {
		sqlDB.Close()
	}
}
