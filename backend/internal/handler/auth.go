package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles all authentication-related endpoints
type AuthHandler struct {
	jwtManager          *utils.JWTManager
	tokenStorageService *service.TokenStorageService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(jwtManager *utils.JWTManager, tokenStorageService *service.TokenStorageService) *AuthHandler {
	return &AuthHandler{
		jwtManager:          jwtManager,
		tokenStorageService: tokenStorageService,
	}
}

// demoUsers is a hardcoded user store for demo purposes.
// Replace with a real DB lookup once a user repository is ready.
var demoUsers = map[string]*models.User{
	"admin@demo.com": {
		ID:           "80882efa-b305-42ad-ad5a-a19cda59590a",
		Email:        "admin@demo.com",
		PasswordHash: mustHash("password123"),
		Role:         models.RoleAdmin,
	},
	"user@demo.com": {
		ID:           "7ca3bfd8-6d00-43f3-a463-cc57cc872898",
		Email:        "user@demo.com",
		PasswordHash: mustHash("password123"),
		Role:         models.RoleUser,
	},
}

func mustHash(pw string) string {
	hash, err := utils.HashPassword(pw)
	if err != nil {
		panic(err)
	}
	return hash
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email + password, receive access & refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Credentials"
// @Success      200   {object}  models.LoginResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	// Lookup user (demo store — swap with DB query later)
	user, ok := demoUsers[req.Email]
	if !ok {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid email or password"))
		return
	}

	// Verify password
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid email or password"))
		return
	}

	sessionID := uuid.New().String()
	userInfo := utils.UserInfo{
		UserID:        user.ID,
		Email:         user.Email,
		WalletAddress: user.WalletAddress,
		Role:          user.Role,
	}

	// Generate access token
	accessToken, err := h.jwtManager.GenerateToken(userInfo, sessionID)
	if err != nil {
		log.Printf("failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to generate token"))
		return
	}

	// Generate refresh token
	refreshToken, err := h.jwtManager.GenerateRefreshToken(userInfo, sessionID)
	if err != nil {
		log.Printf("failed to generate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to generate refresh token"))
		return
	}

	// Store access token JTI in Redis (so it can be validated / revoked)
	accessClaims, _ := h.jwtManager.ValidateToken(accessToken)
	if accessClaims != nil {
		ttl := time.Duration(h.jwtManager.AccessExpirationSeconds) * time.Second
		if err := h.tokenStorageService.StoreJTI(accessClaims.ID, ttl); err != nil {
			log.Printf("failed to store access JTI: %v", err)
		}
	}

	// Store refresh token JTI in Redis & Create PostgreSQL Session
	refreshClaims, _ := h.jwtManager.ValidateToken(refreshToken)
	if refreshClaims != nil {
		ttl := time.Duration(h.jwtManager.RefreshExpirationSeconds) * time.Second
		ipAddr := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		session := &models.Session{
			SessionID: sessionID,
			UserID:    user.ID,
			DeviceID:  nil, // can be added to LoginRequest if needed
			JTI:       refreshClaims.ID,
			IPAddress: &ipAddr,
			UserAgent: &userAgent,
			ExpiresAt: time.Now().Add(ttl),
		}

		if err := h.tokenStorageService.StoreSession(c.Request.Context(), session, ttl); err != nil {
			log.Printf("failed to store refresh JTI & session: %v", err)
		}
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    h.jwtManager.AccessExpirationSeconds,
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Exchange a valid refresh token for a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RefreshRequest  true  "Refresh token"
// @Success      200   {object}  models.RefreshResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      401   {object}  models.ErrorResponse
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err.Error()))
		return
	}

	// Validate the refresh token
	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid or expired refresh token"))
		return
	}

	if claims.TokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Token is not a refresh token"))
		return
	}

	// Check refresh JTI is still valid (not revoked)
	if !h.tokenStorageService.IsValid(claims.ID) {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Refresh token has been revoked"))
		return
	}

	userInfo := utils.UserInfo{
		UserID:        claims.UserID,
		Email:         claims.Email,
		WalletAddress: claims.WalletAddress,
		Role:          claims.Role,
	}

	// Generate new access token with same session
	newAccessToken, err := h.jwtManager.GenerateToken(userInfo, claims.SessionID)
	if err != nil {
		log.Printf("failed to generate new access token: %v", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Failed to generate token"))
		return
	}

	// Store new access token JTI
	newClaims, _ := h.jwtManager.ValidateToken(newAccessToken)
	if newClaims != nil {
		ttl := time.Duration(h.jwtManager.AccessExpirationSeconds) * time.Second
		if err := h.tokenStorageService.StoreJTI(newClaims.ID, ttl); err != nil {
			log.Printf("failed to store new access JTI: %v", err)
		}
	}

	c.JSON(http.StatusOK, models.RefreshResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   h.jwtManager.AccessExpirationSeconds,
	})
}

// Logout godoc
// @Summary      Logout
// @Description  Revoke the current access token (requires Authorization: Bearer <token>)
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.MessageResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract token from header directly (middleware hasn't run on this route)
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader || tokenString == "" {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Authorization header required"))
		return
	}

	claims, err := h.jwtManager.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Invalid token"))
		return
	}
	if claims.TokenType != "access" {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Must provide an access token"))
		return
	}

	// Revoke the Access Token JTI
	if err := h.tokenStorageService.RevokeJTI(claims.ID); err != nil {
		log.Printf("failed to revoke access JTI: %v", err)
	}

	// Revoke the entire Session (which also removes Refresh Token JTI)
	if claims.SessionID != "" {
		if err := h.tokenStorageService.RevokeSessionBySessionID(c.Request.Context(), claims.SessionID); err != nil {
			log.Printf("failed to revoke session %s: %v", claims.SessionID, err)
		}
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Logged out successfully"})
}
