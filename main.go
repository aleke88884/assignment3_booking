package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"smartbooking/config"
	"smartbooking/internal/database"
	"smartbooking/internal/handler"
	"smartbooking/internal/repository"
	"smartbooking/internal/service"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title SmartBooking API
// @version 1.0
// @description Online booking management system for small businesses
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@smartbooking.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.New(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Printf("Connected to PostgreSQL database at %s:%d", cfg.Database.Host, cfg.Database.Port)

	// Initialize repositories with database connection
	userRepo := repository.NewUserRepository(db.DB)
	resourceRepo := repository.NewResourceRepository(db.DB)
	bookingRepo := repository.NewBookingRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	resourceService := service.NewResourceService(resourceRepo)
	bookingService := service.NewBookingService(bookingRepo, resourceRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	resourceHandler := handler.NewResourceHandler(resourceService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// Setup router
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// User routes
	mux.HandleFunc("GET /api/users", userHandler.List)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetByID)
	mux.HandleFunc("GET /api/users/{id}/bookings", bookingHandler.ListByUser)

	// Resource routes
	mux.HandleFunc("GET /api/resources", resourceHandler.List)
	mux.HandleFunc("POST /api/resources", resourceHandler.Create)
	mux.HandleFunc("GET /api/resources/{id}", resourceHandler.GetByID)
	mux.HandleFunc("DELETE /api/resources/{id}", resourceHandler.Delete)

	// Booking routes
	mux.HandleFunc("GET /api/bookings", bookingHandler.ListAll)
	mux.HandleFunc("POST /api/bookings", bookingHandler.Create)
	mux.HandleFunc("GET /api/bookings/{id}", bookingHandler.GetByID)
	mux.HandleFunc("POST /api/bookings/{id}/cancel", bookingHandler.Cancel)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Swagger documentation
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// Start background worker for statistics tracking (demonstrates concurrency)
	go startStatisticsWorker(bookingService, resourceService, userService)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("SmartBooking server starting on %s", addr)
	log.Printf("Available endpoints:")
	log.Printf("  POST /api/auth/register - Register new user")
	log.Printf("  POST /api/auth/login    - User login")
	log.Printf("  GET  /api/users         - List all users")
	log.Printf("  GET  /api/users/{id}    - Get user by ID")
	log.Printf("  GET  /api/resources     - List all resources")
	log.Printf("  POST /api/resources     - Create resource")
	log.Printf("  GET  /api/bookings      - List all bookings")
	log.Printf("  POST /api/bookings      - Create booking")
	log.Printf("  GET  /health            - Health check")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// startStatisticsWorker runs in the background and periodically logs system statistics
// This demonstrates the use of goroutines for concurrent background processing
func startStatisticsWorker(bookingService service.BookingService, resourceService service.ResourceService, userService service.UserService) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Printf("Background statistics worker started")

	for range ticker.C {
		ctx := context.Background()

		// Collect statistics
		bookings, _ := bookingService.ListAll(ctx)
		resources, _ := resourceService.List(ctx)
		users, _ := userService.List(ctx)

		// Count active bookings
		activeBookings := 0
		for _, booking := range bookings {
			if booking.Status != "cancelled" {
				activeBookings++
			}
		}

		// Log statistics
		log.Printf("[STATS] Total Users: %d, Total Resources: %d, Total Bookings: %d, Active Bookings: %d",
			len(users), len(resources), len(bookings), activeBookings)
	}
}
