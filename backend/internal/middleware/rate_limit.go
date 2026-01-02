package middleware

import (
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(cfg *config.RateLimit) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiter := rate.NewLimiter(
		rate.Limit(cfg.RPS),
		cfg.Burst,
	)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
