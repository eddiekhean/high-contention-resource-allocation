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

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/client"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/handler"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/middleware"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service/maze"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/storage"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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

	//ConnectRedis
	rdb, err := client.NewRedisClient(&cfg.RedisConfig)
	store := storage.NewSlotStore(rdb)

	// Services
	simulateService := service.NewSimulateService(logger, store)
	mazeService := maze.NewMazeService(cfg, logger)

	// Handlers
	simulateHandler := handler.NewSimulateHandler(simulateService, logger)
	mazeHandler := handler.NewMazeHandler(mazeService, logger)

	if err != nil {
		logger.Fatalf("redis connect failed: %v", err)
	}

	if rdb != nil {
		logger.Info("redis connected")
	}

	r := gin.New()

	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.RateLimitMiddleware(&cfg.RateLimit),
		middleware.CORSMiddleware(cfg.Cors.AllowedOrigins),
	)
	r.GET("/health", handler.HealthCheck)

	public := r.Group("/api/v1/public")
	{
		simulate := public.Group("/simulate")
		{
			simulate.POST("/run", simulateHandler.Simulate)
		}
		leetcode := public.Group("/leetcode")
		{
			maze := leetcode.Group("/maze")
			{
				maze.POST("/submit", mazeHandler.Submit)
				maze.POST("/generate", mazeHandler.Generate)
			}
		}
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exiting")
}
