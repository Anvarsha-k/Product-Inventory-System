package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Anvarsha-k/Product-Inventory-System/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN not set")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	
	if err := db.AutoMigrate(&models.Product{}, &models.Variant{}, &models.VariantOption{}, &models.SubVariant{}, &models.StockTransaction{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
	DB = db
	fmt.Println("DB connected")
}
