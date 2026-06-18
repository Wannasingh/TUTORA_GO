package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/Wannasingh/TUTORA_GO/backend/config"
	delivery "github.com/Wannasingh/TUTORA_GO/backend/delivery/http"
	"github.com/Wannasingh/TUTORA_GO/backend/repository"
	"github.com/Wannasingh/TUTORA_GO/backend/usecase"
)

func main() {
	// 1. Load Configurations
	cfg := config.LoadConfig()

	// 2. Setup DB Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Connecting to database at %s...", cfg.DBConnString)
	dbConn, err := pgx.Connect(ctx, cfg.DBConnString)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v. Running in mockup mode or offline.", err)
	} else {
		defer dbConn.Close(context.Background())
		log.Println("Database connection established successfully!")
	}

	// 3. Initialize layers manually (Dependency Injection)
	userRepository := repository.NewPostgresUserRepository(dbConn)
	tutorRepository := repository.NewPostgresTutorRepository(dbConn)

	userUsecase := usecase.NewUserUsecase(userRepository)
	tutorUsecase := usecase.NewTutorUsecase(tutorRepository, userRepository)
	authUsecase := usecase.NewAuthUsecase(userRepository, cfg)

	// 4. Setup Web Server (Gin)
	r := gin.Default()

	// Simple Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
			"schema": cfg.DBSchema,
		})
	})

	// 5. Register Delivery HTTP handlers
	delivery.NewHttpHandler(r, userUsecase, tutorUsecase, authUsecase)

	// 6. Start Server
	log.Printf("Server is running on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Since DB may not always be available during initial build or tests, let's keep mock definitions in separate files or code to avoid build errors.
// Wait, to compile fine even if DB connection is unavailable at run time, we just connect normally.
