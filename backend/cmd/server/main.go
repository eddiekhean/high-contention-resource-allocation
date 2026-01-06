package main

import (
	"flag"
	_ "fmt"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/client"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/handler"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/middleware"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
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
	simulateService := service.NewSimulateService(logger, store)
	simulateHandler := handler.NewSimulateHandler(simulateService, logger)
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
		middleware.CORSMiddleware(),
	)
	r.GET("/health", handler.HealthCheck)
	r.POST("/simulate", simulateHandler.Simulate)
	r.Run(":8080")
}
