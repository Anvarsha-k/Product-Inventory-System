package main

import (
	"log"
	"os"

	"github.com/Anvarsha-k/Product-Inventory-System/internal/db"
	"github.com/Anvarsha-k/Product-Inventory-System/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

		err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found or unable to load it")
	}
	db.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := fiber.New()

	api := app.Group("/api")
	api.Post("/products", handlers.CreateProduct)
	api.Get("/products", handlers.ListProducts)
	api.Post("/stock/adjust", handlers.AdjustStock)

	log.Printf("Server running on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
