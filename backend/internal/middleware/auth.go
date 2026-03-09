package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// TokenValidationMiddleware validates JWT tokens using only JTI lookup (very fast)
func TokenValidationMiddleware(jwtManager *utils.JWTManager, tokenStorageService *service.TokenStorageService, roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header required"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Bearer token required"))
			c.Abort()
			return
		}

		// Parse and validate JWT signature first
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			log.Printf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token"))
			c.Abort()
			return
		}
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token type"))
			c.Abort()
			return
		}

		// Extract JTI from claims
		jti := claims.ID // JTI is stored in RegisteredClaims.ID
		if jti == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Token missing JTI"))
			c.Abort()
			return
		}

		// Fast JTI validation - just check if JTI exists and is active
		if !tokenStorageService.IsValid(jti) {
			log.Printf("JTI validation failed: %s", jti)
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Token has been revoked"))
			c.Abort()
			return
		}

		// Add user info to context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("session_id", claims.SessionID)
		c.Set("user_role", claims.Role)
		c.Set("wallet_address", claims.WalletAddress)
		c.Set("email", claims.Email)
		if len(roles) > 0 && !containsRole(roles, claims.Role) {
			c.JSON(http.StatusForbidden, models.NewErrorResponse("Insufficient permissions"))
			c.Abort()
			return
		}

		c.Set("jti", jti)
		c.Next()
	}
}
func containsRole(roles []models.UserRole, role models.UserRole) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
