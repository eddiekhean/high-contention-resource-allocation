package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/client"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/db"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/handler"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/middleware"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/repository"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/storage"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found or error loading it")
	}

	var (
		configFile = flag.String("config", "config.yaml", "Path to YAML configuration file")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadFromFile(*configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	logger := utils.SetupLogger(cfg)
	logger.Info("Configuration loaded successfully")

	// Redis
	rdb, err := client.NewRedisClient(&cfg.RedisConfig)
	if err != nil {
		logger.Fatalf("redis connect failed: %v", err)
	}
	store := storage.NewSlotStore(rdb)
	simulateService := service.NewSimulateService(logger, store)
	simulateHandler := handler.NewSimulateHandler(simulateService, logger)

	if rdb != nil {
		logger.Info("redis connected")
	}

	// S3
	s3Client, err := client.NewS3Client(&cfg.S3)
	if err != nil {
		logger.Fatalf("s3 client initialization failed: %v", err)
	}
	if s3Client != nil {
		logger.Info("s3 connected")
	} else {
		logger.Warn("s3 is disabled")
	}

	// Postgres
	pg, err := db.NewPostgres(&cfg.Postgres)
	if err != nil {
		logger.Fatalf("postgres client initialization failed: %v", err)
	}
	if pg != nil {
		logger.Info("postgres connected")
		if err := db.RunMigration(pg); err != nil {
			logger.Fatalf("db migration failed: %v", err)
		}
	} else {
		logger.Warn("postgres is disabled")
	}
	// Repository
	imageRepository := repository.NewPgImageRepository(pg)
	imageService := service.NewImageService(imageRepository, &cfg.ImageConfig, s3Client)
	mazeHandler := handler.NewMazeHandler(imageService, logger)

	// Gin
	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.RateLimitMiddleware(&cfg.RateLimit),
		middleware.CORSMiddleware(),
	)

	// system
	r.GET("/health", handler.HealthCheck)

	// simulate (generic playground)
	r.POST("/simulate", simulateHandler.Simulate)

	// Images

	// ===== LEETCODE PAGE =====
	leetcode := r.Group("/leetcode")
	{
		maze := leetcode.Group("/maze")
		{
			images := maze.Group("/images")
			{
				images.POST("/match", mazeHandler.MatchImage)
				images.POST("/upload-url", mazeHandler.GetUploadURL)
				images.POST("", mazeHandler.CommitImage)
			}
		}
	}
	// HTTP server (QUAN TRỌNG)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run server async
	go func() {
		logger.Info("HTTP server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen error: %v", err)
		}
	}()

	// ===== GRACEFUL SHUTDOWN =====
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown failed: %v", err)
	}

	// Close DB
	if pg != nil {
		pg.Close()
		logger.Info("postgres closed")
	}

	// Close Redis
	if rdb != nil {
		_ = rdb.Close()
		logger.Info("redis closed")
	}

	logger.Info("shutdown complete")
}
