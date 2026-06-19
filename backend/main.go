package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Wannasingh/TUTORA_GO/backend/config"
	delivery "github.com/Wannasingh/TUTORA_GO/backend/delivery/http"
	"github.com/Wannasingh/TUTORA_GO/backend/repository"
	"github.com/Wannasingh/TUTORA_GO/backend/usecase"
	"github.com/Wannasingh/TUTORA_GO/backend/utils"
)

func main() {
	// Disable AWS chunked encoding checks which are unsupported by OCI Object Storage
	os.Setenv("AWS_REQUEST_CHECKSUM_CALCULATION", "when_required")
	os.Setenv("AWS_RESPONSE_CHECKSUM_VALIDATION", "when_required")

	// 1. Load Configurations
	cfg := config.LoadConfig()

	// 2. Setup DB Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Connecting to database pool at %s...", cfg.DBConnString)
	dbPool, err := pgxpool.New(ctx, cfg.DBConnString)
	if err != nil {
		log.Printf("Warning: Failed to create connection pool: %v. Running in mockup mode or offline.", err)
	} else {
		defer dbPool.Close()
		log.Println("Database connection pool established successfully!")
	}

	// 3. Initialize layers manually (Dependency Injection)
	userRepository := repository.NewPostgresUserRepository(dbPool)
	tutorRepository := repository.NewPostgresTutorRepository(dbPool)
	postRepository := repository.NewPostgresPostRepository(dbPool)

	userUsecase := usecase.NewUserUsecase(userRepository)
	tutorUsecase := usecase.NewTutorUsecase(tutorRepository, userRepository)
	authUsecase := usecase.NewAuthUsecase(userRepository, cfg)
	postUsecase := usecase.NewPostUsecase(postRepository)

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

	// 4.5. Initialize OCI storage service
	storageService, err := utils.NewOCIStorageService(
		cfg.OCIS3AccessKeyID,
		cfg.OCIS3SecretAccessKey,
		cfg.OCIS3Region,
		cfg.OCIS3BucketName,
		cfg.OCIS3Endpoint,
	)
	if err != nil {
		log.Fatalf("Failed to initialize OCI storage service: %v", err)
	}

	// 5. Register Delivery HTTP handlers
	delivery.NewHttpHandler(r, userUsecase, tutorUsecase, authUsecase, postUsecase, storageService, cfg)

	// 6. Start Server
	log.Printf("Server is running on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Since DB may not always be available during initial build or tests, let's keep mock definitions in separate files or code to avoid build errors.
// Wait, to compile fine even if DB connection is unavailable at run time, we just connect normally.
