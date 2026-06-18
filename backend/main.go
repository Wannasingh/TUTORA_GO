package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/haru/bytestutor/backend/config"
	delivery "github.com/haru/bytestutor/backend/delivery/http"
	"github.com/haru/bytestutor/backend/repository"
	"github.com/haru/bytestutor/backend/usecase"
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
	delivery.NewHttpHandler(r, userUsecase, tutorUsecase)

	// 6. Start Server
	log.Printf("Server is running on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Since DB may not always be available during initial build or tests, let's keep mock definitions in separate files or code to avoid build errors.
// Wait, to compile fine even if DB connection is unavailable at run time, we just connect normally.
