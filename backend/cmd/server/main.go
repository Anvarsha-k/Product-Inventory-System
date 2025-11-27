package main

import (
	"log"
	"os"

	"github.com/Anvarsha-k/Product-Inventory-System/internal/db"
	"github.com/Anvarsha-k/Product-Inventory-System/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := fiber.New()

	api := app.Group("/api")
	api.Post("/products", handlers.CreateProduct)

	log.Printf("Server running on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
