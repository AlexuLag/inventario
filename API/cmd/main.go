package main

import (
	"database/sql"
	"log"
	"net/http"

	"inventario/internal/infrastructure/repository"
	"inventario/internal/interface/handler"
	"inventario/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database
	db, err := sql.Open("sqlite3", "products.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables
	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	productRepo, err := repository.NewSQLiteProductRepository("products.db")
	if err != nil {
		log.Fatal(err)
	}
	defer productRepo.Close()

	userRepo := repository.NewSQLiteUserRepository(db)

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase)
	userHandler := handler.NewUserHandler(userUseCase)

	// Initialize router
	r := chi.NewRouter()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(r chi.Router) {
		// Product routes
		r.Route("/products", func(r chi.Router) {
			r.Get("/", productHandler.GetAllProducts)
			r.Post("/", productHandler.CreateProduct)
			r.Get("/{id}", productHandler.GetProduct)
			r.Put("/{id}", productHandler.UpdateProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
		})

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetAllUsers)
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})
	})

	// Start server
	log.Println("Starting server...")
	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func createTables(db *sql.DB) error {
	// Create products table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			code TEXT NOT NULL UNIQUE,
			image_url TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			role TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	return err
}
