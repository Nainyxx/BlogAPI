package db

import (
	db "BlogAPI/init"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	cfg := db.LoadDBCfg()

	// Подключаемся
	database, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = database
	log.Println("✅ Database connected successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
