package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"smartbooking/config"
	_ "smartbooking/docs"
	"smartbooking/internal/database"
	"smartbooking/internal/handler"
	"smartbooking/internal/repository"
	"smartbooking/internal/service"
	"smartbooking/internal/storage"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	cfg := config.Load()

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

	// Проверяем схему БД
	ctx := context.Background()
	if err := db.VerifySchema(ctx); err != nil {
		log.Printf("WARNING: Schema verification failed: %v", err)
		log.Printf("Убедитесь что миграции выполнены (они должны запуститься автоматически при первом запуске PostgreSQL)")
	}

	// Выводим статистику БД
	stats, err := db.GetDatabaseStats(ctx)
	if err == nil {
		log.Printf("БД статистика: Users=%d, Resources=%d, Bookings=%d",
			stats.UsersCount, stats.ResourcesCount, stats.BookingsCount)
	}

	// Инициализируем storage (S3/MinIO)
	var storageService storage.StorageService
	if cfg.Storage.Type == "local" {
		log.Println("Using local storage")
		storageService = storage.NewLocalStorage("/uploads", "http://localhost:8080/uploads")
	} else {
		log.Printf("Connecting to %s storage at %s", cfg.Storage.Type, cfg.Storage.Endpoint)
		storageService, err = storage.NewS3Storage(storage.S3Config{
			Endpoint:        cfg.Storage.Endpoint,
			AccessKeyID:     cfg.Storage.AccessKeyID,
			SecretAccessKey: cfg.Storage.SecretAccessKey,
			BucketName:      cfg.Storage.BucketName,
			Region:          cfg.Storage.Region,
			UseSSL:          cfg.Storage.UseSSL,
			PublicURL:       cfg.Storage.PublicURL,
		})
		if err != nil {
			log.Fatalf("Failed to initialize storage: %v", err)
		}
		log.Printf("✓ Storage initialized successfully")
	}

	userRepo := repository.NewUserRepository(db.DB)
	resourceRepo := repository.NewResourceRepository(db.DB)
	bookingRepo := repository.NewBookingRepository(db.DB)
	photoRepo := repository.NewPhotoRepository(db.DB)
	reviewRepo := repository.NewReviewRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	resourceService := service.NewResourceService(resourceRepo)
	bookingService := service.NewBookingService(bookingRepo, resourceRepo)
	photoService := service.NewPhotoService(photoRepo, storageService)
	reviewService := service.NewReviewService(reviewRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	resourceHandler := handler.NewResourceHandler(resourceService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	photoHandler := handler.NewPhotoHandler(photoService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	mux := http.NewServeMux()

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	mux.HandleFunc("GET /api/users", userHandler.List)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetByID)
	mux.HandleFunc("GET /api/users/{id}/bookings", bookingHandler.ListByUser)

	mux.HandleFunc("GET /api/resources", resourceHandler.List)
	mux.HandleFunc("POST /api/resources", resourceHandler.Create)
	mux.HandleFunc("GET /api/resources/{id}", resourceHandler.GetByID)
	mux.HandleFunc("DELETE /api/resources/{id}", resourceHandler.Delete)

	mux.HandleFunc("GET /api/bookings", bookingHandler.ListAll)
	mux.HandleFunc("POST /api/bookings", bookingHandler.Create)
	mux.HandleFunc("GET /api/bookings/{id}", bookingHandler.GetByID)
	mux.HandleFunc("POST /api/bookings/{id}/cancel", bookingHandler.Cancel)

	mux.HandleFunc("POST /api/photos/upload", photoHandler.UploadPhoto)
	mux.HandleFunc("GET /api/resources/{resource_id}/photos", photoHandler.GetResourcePhotos)
	mux.HandleFunc("DELETE /api/photos/{id}", photoHandler.DeletePhoto)
	mux.HandleFunc("PUT /api/photos/{id}/primary", photoHandler.SetPrimaryPhoto)

	mux.HandleFunc("GET /api/reviews", reviewHandler.GetByResource)
	mux.HandleFunc("POST /api/reviews", reviewHandler.Create)
	mux.HandleFunc("GET /api/reviews/{id}", reviewHandler.GetByID)
	mux.HandleFunc("PUT /api/reviews/{id}", reviewHandler.Update)
	mux.HandleFunc("DELETE /api/reviews/{id}", reviewHandler.Delete)
	mux.HandleFunc("GET /api/resources/{resource_id}/reviews", reviewHandler.GetByResource)
	mux.HandleFunc("GET /api/resources/{resource_id}/rating", reviewHandler.GetResourceAverageRating)
	mux.HandleFunc("GET /api/users/{user_id}/reviews", reviewHandler.GetByUser)

	mux.HandleFunc("GET /api/categories", categoryHandler.List)
	mux.HandleFunc("POST /api/categories", categoryHandler.Create)
	mux.HandleFunc("GET /api/categories/{id}", categoryHandler.GetByID)
	mux.HandleFunc("PUT /api/categories/{id}", categoryHandler.Update)
	mux.HandleFunc("DELETE /api/categories/{id}", categoryHandler.Delete)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	go startStatisticsWorker(bookingService, resourceService, userService)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("SmartBooking server starting on %s", addr)
	log.Printf("Available endpoints:")
	log.Printf("  POST /api/auth/register              - Register new user")
	log.Printf("  POST /api/auth/login                 - User login")
	log.Printf("  GET  /api/users                      - List all users")
	log.Printf("  GET  /api/users/{id}                 - Get user by ID")
	log.Printf("  GET  /api/resources                  - List all resources")
	log.Printf("  POST /api/resources                  - Create resource")
	log.Printf("  GET  /api/bookings                   - List all bookings")
	log.Printf("  POST /api/bookings                   - Create booking")
	log.Printf("  POST /api/photos/upload              - Upload photo")
	log.Printf("  GET  /api/resources/{id}/photos      - Get resource photos")
	log.Printf("  DELETE /api/photos/{id}              - Delete photo")
	log.Printf("  GET  /health                         - Health check")
	log.Printf("  GET  /swagger/                       - API documentation")
	if cfg.Storage.Type == "minio" {
		log.Printf("  MinIO Console: http://localhost:9001 (admin/admin)")
	}

	if err := http.ListenAndServe(addr, corsMiddleware(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func startStatisticsWorker(bookingService service.BookingService, resourceService service.ResourceService, userService service.UserService) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Printf("Background statistics worker started")

	for range ticker.C {
		ctx := context.Background()

		bookings, _ := bookingService.ListAll(ctx)
		resources, _ := resourceService.List(ctx)
		users, _ := userService.List(ctx)

		activeBookings := 0
		for _, booking := range bookings {
			if booking.Status != "cancelled" {
				activeBookings++
			}
		}

		log.Printf("[STATS] Total Users: %d, Total Resources: %d, Total Bookings: %d, Active Bookings: %d",
			len(users), len(resources), len(bookings), activeBookings)
	}
}
