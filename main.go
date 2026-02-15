package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"

	"wearlab_backend/internal/handler"
	"wearlab_backend/internal/repository"
	"wearlab_backend/internal/service"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "wearlab"
	password = "wearlabbro30102001"
	dbname   = "wearlabdatabase"
)

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	commonRepo := repository.NewCommonRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo)
	commonService := service.NewCommonService(commonRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	commonHandler := handler.NewCommonHandler(commonService)

	app := fiber.New()

	app.Use(cors.New())

	// Product routes
	app.Get("/product", productHandler.GetProducts)
	app.Post("/product", productHandler.CreateProduct)
	app.Post("/product/sell", productHandler.SellProduct)
	app.Put("/product/bulk-price-update", productHandler.BulkPriceUpdate)

	// Common routes
	app.Get("/type", commonHandler.GetTypes)

	// Start server
	app.Listen(":8080")
}
