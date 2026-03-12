package main

import (
	"context"
	"flag"
	_ "fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/eddiekhean/high-contention-resource-allocation-backend/docs" // Import swagger docs

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/client"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/database"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/handler"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/middleware"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service/maze"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/storage"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           High Contention Resource Allocation API
// @version         1.0
// @description     API Server for High Contention Resource Allocation app.
// @host      localhost:9091
// @BasePath  /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	var (
		configFile = flag.String("config", "config.yaml", "Path to YAML configuration file (default: config.yaml)")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadFromFile(*configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger based on configuration
	logger := utils.SetupLogger(cfg)
	logger.Info("Configuration loaded successfully")

	// ConnectRedis
	rdb, err := client.NewRedisClient(&cfg.RedisConfig)
	if err != nil {
		logger.Warnf("redis connect failed (falling back to in-memory token store): %v", err)
		rdb = nil
	}
	if rdb != nil {
		logger.Info("redis connected")
	} else {
		logger.Warn("running WITHOUT Redis — token storage is in-memory only (not suitable for production)")
	}

	store := storage.NewSlotStore(rdb)

	// ── PostgreSQL ────────────────────────────────────────────────────────────
	pgPool, err := client.NewPostgresPool(&cfg.Postgres)
	if err != nil {
		logger.Warnf("postgres connect failed: %v — running without PostgreSQL", err)
	} else if pgPool != nil {
		logger.Info("postgres connected")
		// Auto-run SQL migrations on startup (idempotent, skips already applied)
		if err := database.MigrateUp(pgPool); err != nil {
			logger.Fatalf("postgres migrations failed: %v", err)
		}
		logger.Info("postgres migrations applied successfully")
	}

	// ── MongoDB ───────────────────────────────────────────────────────────────
	var messageStore *storage.MessageStore
	mongoClient, err := client.NewMongoClient(&cfg.Mongo)
	if err != nil {
		logger.Warnf("mongo connect failed: %v — running without MongoDB", err)
	} else if mongoClient != nil {
		logger.Info("mongo connected")
		mongoDB := mongoClient.Database(cfg.Mongo.Database)
		messageStore = storage.NewMessageStore(mongoDB)
		// Ensure indexes are created (idempotent)
		ctxIdx, cancelIdx := context.WithTimeout(context.Background(), 10*time.Second)
		if err := messageStore.EnsureIndexes(ctxIdx); err != nil {
			logger.Warnf("mongo EnsureIndexes failed: %v", err)
		} else {
			logger.Info("mongo indexes ensured")
		}
		cancelIdx()
	}
	_ = messageStore // will be used by handlers when implemented

	// Initialize JWT Manager from config
	jwtManager, err := utils.NewJWTManagerFromConfig(
		cfg.JWT.SigningMethod,
		"", // HMAC secret (not used for RSA)
		cfg.JWT.PrivateKeyPath,
		cfg.JWT.PublicKeyPath,
		"", // EC private key path (not used for RSA)
		"", // EC public key path (not used for RSA)
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)
	if err != nil {
		logger.Fatalf("Failed to initialize JWT manager: %v", err)
	}
	logger.Info("JWT manager initialized successfully")

	// Initialize Token Storage Service (Redis-backed + PostgreSQL)
	tokenStorageService := service.NewTokenStorageService(rdb, pgPool)

	// Services
	simulateService := service.NewSimulateService(logger, store)
	mazeService := maze.NewMazeService(cfg, logger)

	// Handlers
	simulateHandler := handler.NewSimulateHandler(simulateService, logger)
	mazeHandler := handler.NewMazeHandler(mazeService, logger)
	authHandler := handler.NewAuthHandler(jwtManager, tokenStorageService)

	r := gin.New()

	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.RateLimitMiddleware(&cfg.RateLimit),
		middleware.CORSMiddleware(cfg.Cors.AllowedOrigins),
	)

	// Health check (public)
	r.GET("/health", handler.HealthCheck)

	// Swagger documentation (public)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes (public — no JWT required)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}

	// Public API routes (no auth)
	public := r.Group("/api/v1/public")
	{
		simulate := public.Group("/simulate")
		{
			simulate.POST("/run", simulateHandler.Simulate)
		}
		leetcode := public.Group("/leetcode")
		{
			mazeGroup := leetcode.Group("/maze")
			{
				mazeGroup.POST("/submit", mazeHandler.Submit)
				mazeGroup.POST("/generate", mazeHandler.Generate)
			}
		}
	}

	// Protected routes — require valid JWT (any role)
	protected := r.Group("/api/v1/protected")
	protected.Use(middleware.TokenValidationMiddleware(jwtManager, tokenStorageService))
	{
		protected.GET("/health", handler.HealthCheck)
	}

	// Admin-only routes — require valid JWT with admin role
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.TokenValidationMiddleware(jwtManager, tokenStorageService, models.RoleAdmin))
	{
		_ = admin // add admin-specific routes here
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine to enable graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	logger.Infof("Server started on :%s", port)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exiting")
}
