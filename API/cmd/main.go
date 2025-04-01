package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"inventario/internal/infrastructure/repository"
	"inventario/internal/interface/handler"
	"inventario/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize MySQL connection
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	// Initialize repositories
	productRepo := repository.NewMySQLProductRepository(db)
	userRepo := repository.NewMySQLUserRepository(db)
	stockRepo := repository.NewMySQLStockRepository(db)
	providerRepo := repository.NewMySQLProviderRepository(db)

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	stockUseCase := usecase.NewStockUseCase(stockRepo)
	providerUseCase := usecase.NewProviderUseCase(providerRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	stockHandler := handler.NewStockHandler(stockUseCase)
	providerHandler := handler.NewProviderHandler(providerUseCase)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Route("/api", func(r chi.Router) {
		// Product routes
		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHandler.CreateProduct)
			r.Get("/", productHandler.GetAllProducts)
			r.Get("/{id}", productHandler.GetProduct)
			r.Put("/{id}", productHandler.UpdateProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
		})

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.GetAllUsers)
			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		// Stock routes
		r.Route("/stocks", func(r chi.Router) {
			r.Post("/", stockHandler.CreateStock)
			r.Get("/", stockHandler.GetAllStocks)
			r.Get("/{id}", stockHandler.GetStock)
			r.Put("/{id}", stockHandler.UpdateStock)
			r.Delete("/{id}", stockHandler.DeleteStock)
			r.Get("/product/{productId}", stockHandler.GetStocksByProductID)
			r.Get("/serial/{serial}", stockHandler.GetStockBySerial)
		})

		// Provider routes
		r.Route("/providers", func(r chi.Router) {
			r.Post("/", providerHandler.CreateProvider)
			r.Get("/", providerHandler.GetAllProviders)
			r.Get("/{id}", providerHandler.GetProvider)
			r.Put("/{id}", providerHandler.UpdateProvider)
			r.Delete("/{id}", providerHandler.DeleteProvider)
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
