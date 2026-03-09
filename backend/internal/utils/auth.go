package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	UserID        string          `json:"user_id"`
	Email         string          `json:"email"`
	WalletAddress string          `json:"wallet_address,omitempty"`
	Role          models.UserRole `json:"role,omitempty"` // User role (admin, user)
}

type JWTClaims struct {
	UserInfo
	TokenType string `json:"token_type,omitempty"` // Optional for refresh tokens
	SessionID string `json:"session_id"`           // Add SessionID for session management

	jwt.RegisteredClaims
}

// JWTManager handles JWT operations with different signing methods
type JWTManager struct {
	SigningMethod            string
	Secret                   string
	PrivateKey               interface{} // Can be *rsa.PrivateKey or *ecdsa.PrivateKey
	PublicKey                interface{} // Can be *rsa.PublicKey or *ecdsa.PublicKey
	AccessExpirationSeconds  int
	RefreshExpirationSeconds int
}

// NewJWTManager creates a new JWT manager with HMAC signing
func NewJWTManagerHMAC(signingMethod, secret string, accessExpirationSeconds, refreshExpirationSeconds int) *JWTManager {
	return &JWTManager{
		SigningMethod:            signingMethod,
		Secret:                   secret,
		AccessExpirationSeconds:  accessExpirationSeconds,
		RefreshExpirationSeconds: refreshExpirationSeconds,
	}
}

// NewJWTManagerRSA creates a new JWT manager with RSA signing
func NewJWTManagerRSA(signingMethod, privateKeyPath, publicKeyPath string, accessExpirationSeconds, refreshExpirationSeconds int) (*JWTManager, error) {
	privateKey, err := loadRSAPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load RSA private key: %w", err)
	}

	publicKey, err := loadRSAPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load RSA public key: %w", err)
	}

	return &JWTManager{
		SigningMethod:            signingMethod,
		PrivateKey:               privateKey,
		PublicKey:                publicKey,
		AccessExpirationSeconds:  accessExpirationSeconds,
		RefreshExpirationSeconds: refreshExpirationSeconds,
	}, nil
}

// NewJWTManagerECDSA creates a new JWT manager with ECDSA signing
func NewJWTManagerECDSA(signingMethod, privateKeyPath, publicKeyPath string, accessExpirationSeconds, refreshExpirationSeconds int) (*JWTManager, error) {
	privateKey, err := loadECDSAPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ECDSA private key: %w", err)
	}

	publicKey, err := loadECDSAPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ECDSA public key: %w", err)
	}

	return &JWTManager{
		SigningMethod:            signingMethod,
		PrivateKey:               privateKey,
		PublicKey:                publicKey,
		AccessExpirationSeconds:  accessExpirationSeconds,
		RefreshExpirationSeconds: refreshExpirationSeconds,
	}, nil
}

// ExtractJTI extracts JWT ID from token string without full validation
func (j *JWTManager) ExtractJTI(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.ID, nil
}

// GetTokenInfo returns basic token information
func (j *JWTManager) GetTokenInfo(tokenString string) (userID, jti string, err error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.ID, nil
}

// GenerateToken generates a JWT token using the configured signing method
func (j *JWTManager) GenerateToken(userInfo UserInfo, sessionID string) (string, error) {
	claims := JWTClaims{
		UserInfo: userInfo,

		SessionID: sessionID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "jti-" + uuid.New().String(), // JWT ID for token management
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.AccessExpirationSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "cti-module-3",
			Subject:   userInfo.UserID,
		},
	}

	var token *jwt.Token

	switch {
	case strings.HasPrefix(j.SigningMethod, "HS"): // HMAC methods
		switch j.SigningMethod {
		case "HS256":
			token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		case "HS384":
			token = jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
		case "HS512":
			token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		default:
			return "", fmt.Errorf("unsupported HMAC method: %s", j.SigningMethod)
		}
		return token.SignedString([]byte(j.Secret))

	case strings.HasPrefix(j.SigningMethod, "RS"): // RSA methods
		switch j.SigningMethod {
		case "RS256":
			token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		case "RS384":
			token = jwt.NewWithClaims(jwt.SigningMethodRS384, claims)
		case "RS512":
			token = jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
		default:
			return "", fmt.Errorf("unsupported RSA method: %s", j.SigningMethod)
		}
		return token.SignedString(j.PrivateKey)

	case strings.HasPrefix(j.SigningMethod, "ES"): // ECDSA methods
		switch j.SigningMethod {
		case "ES256":
			token = jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		case "ES384":
			token = jwt.NewWithClaims(jwt.SigningMethodES384, claims)
		case "ES512":
			token = jwt.NewWithClaims(jwt.SigningMethodES512, claims)
		default:
			return "", fmt.Errorf("unsupported ECDSA method: %s", j.SigningMethod)
		}
		return token.SignedString(j.PrivateKey)

	default:
		return "", fmt.Errorf("unsupported signing method: %s", j.SigningMethod)
	}
}

// ValidateToken validates a JWT token using the configured signing method
func (j *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		switch {
		case strings.HasPrefix(j.SigningMethod, "HS"):
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(j.Secret), nil

		case strings.HasPrefix(j.SigningMethod, "RS"):
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.PublicKey, nil

		case strings.HasPrefix(j.SigningMethod, "ES"):
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.PublicKey, nil

		default:
			return nil, fmt.Errorf("unsupported signing method: %s", j.SigningMethod)
		}
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateRefreshToken generates a refresh token
func (j *JWTManager) GenerateRefreshToken(user UserInfo, sessionID string) (string, error) {
	// Note: We only store UserID in refresh token for security
	// When refreshing, we'll fetch user details from database
	claims := JWTClaims{
		UserInfo:  user,
		SessionID: sessionID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "jti-" + uuid.New().String(), // JWT ID for token management
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.RefreshExpirationSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "cti-module-3",
			Subject:   user.UserID,
		},
	}

	var token *jwt.Token

	switch {
	case strings.HasPrefix(j.SigningMethod, "HS"):
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Use HS256 for refresh tokens
		return token.SignedString([]byte(j.Secret))

	case strings.HasPrefix(j.SigningMethod, "RS"):
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims) // Use RS256 for refresh tokens
		return token.SignedString(j.PrivateKey)

	case strings.HasPrefix(j.SigningMethod, "ES"):
		token = jwt.NewWithClaims(jwt.SigningMethodES256, claims) // Use ES256 for refresh tokens
		return token.SignedString(j.PrivateKey)

	default:
		return "", fmt.Errorf("unsupported signing method: %s", j.SigningMethod)
	}
}

// NewJWTManagerFromConfig creates a JWTManager from configuration
func NewJWTManagerFromConfig(signingMethod, secret, privateKeyPath, publicKeyPath, ecPrivateKeyPath, ecPublicKeyPath string, accessExpirationSeconds, refreshExpirationSeconds int) (*JWTManager, error) {
	switch {
	case strings.HasPrefix(signingMethod, "HS"):
		return NewJWTManagerHMAC(signingMethod, secret, accessExpirationSeconds, refreshExpirationSeconds), nil

	case strings.HasPrefix(signingMethod, "RS"):
		return NewJWTManagerRSA(signingMethod, privateKeyPath, publicKeyPath, accessExpirationSeconds, refreshExpirationSeconds)

	case strings.HasPrefix(signingMethod, "ES"):
		return NewJWTManagerECDSA(signingMethod, ecPrivateKeyPath, ecPublicKeyPath, accessExpirationSeconds, refreshExpirationSeconds)

	default:
		return nil, fmt.Errorf("unsupported signing method: %s", signingMethod)
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Key loading functions
func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if rsaKey, ok := keyInterface.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, errors.New("not an RSA private key")
	}

	return key, nil
}

func loadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if rsaKey, ok := key.(*rsa.PublicKey); ok {
		return rsaKey, nil
	}

	return nil, errors.New("not an RSA public key")
}

func loadECDSAPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if ecdsaKey, ok := keyInterface.(*ecdsa.PrivateKey); ok {
			return ecdsaKey, nil
		}
		return nil, errors.New("not an ECDSA private key")
	}

	return key, nil
}

func loadECDSAPublicKey(path string) (*ecdsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if ecdsaKey, ok := key.(*ecdsa.PublicKey); ok {
		return ecdsaKey, nil
	}

	return nil, errors.New("not an ECDSA public key")
}
func GenerateRandomState(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
