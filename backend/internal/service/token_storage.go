package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// TokenStorageService handles JTI-based token validation.
// It uses Redis when available, falling back to an in-memory store otherwise.
// It also writes long-lived login sessions to PostgreSQL.
type TokenStorageService struct {
	rdb    *redis.Client
	pgPool *pgxpool.Pool

	// in-memory fallback
	mu       sync.RWMutex
	memStore map[string]time.Time // jti -> expiry
}

func NewTokenStorageService(rdb *redis.Client, pgPool *pgxpool.Pool) *TokenStorageService {
	return &TokenStorageService{
		rdb:      rdb,
		pgPool:   pgPool,
		memStore: make(map[string]time.Time),
	}
}

func jtiKey(jti string) string {
	return fmt.Sprintf("jti:%s", jti)
}

// StoreSession stores a long-lived session in PostgreSQL and its JTI in Redis
func (s *TokenStorageService) StoreSession(ctx context.Context, session *models.Session, ttl time.Duration) error {
	// 1. Store in PostgreSQL for multi-device management and history
	if s.pgPool != nil {
		query := `
			INSERT INTO sessions (session_id, user_id, device_id, jti, ip_address, user_agent, expires_at)
			VALUES ($1, $2, NULLIF($3, '')::uuid, $4, $5, $6, $7)
		`
		var deviceID interface{}
		if session.DeviceID != nil {
			deviceID = *session.DeviceID
		}
		_, err := s.pgPool.Exec(ctx, query,
			session.SessionID, session.UserID, deviceID, session.JTI,
			session.IPAddress, session.UserAgent, session.ExpiresAt)
		if err != nil {
			return fmt.Errorf("failed to store session in db: %w", err)
		}
	}

	// 2. Store JTI in Redis/Mem for fast validation
	return s.StoreJTI(session.JTI, ttl)
}

// StoreJTI stores a JTI with the given TTL in Redis (used for short-lived access tokens)
func (s *TokenStorageService) StoreJTI(jti string, ttl time.Duration) error {
	if s.rdb != nil {
		ctx := context.Background()
		return s.rdb.Set(ctx, jtiKey(jti), "1", ttl).Err()
	}
	// In-memory fallback
	s.mu.Lock()
	defer s.mu.Unlock()
	s.memStore[jti] = time.Now().Add(ttl)
	return nil
}

// IsValid checks whether a JTI is still valid (not revoked, not expired)
func (s *TokenStorageService) IsValid(jti string) bool {
	if s.rdb != nil {
		ctx := context.Background()
		val, err := s.rdb.Exists(ctx, jtiKey(jti)).Result()
		if err != nil {
			return false
		}
		return val > 0
	}
	// In-memory fallback
	s.mu.RLock()
	defer s.mu.RUnlock()
	expiry, ok := s.memStore[jti]
	if !ok {
		return false
	}
	return time.Now().Before(expiry)
}

// RevokeJTI removes a JTI (used for access tokens)
func (s *TokenStorageService) RevokeJTI(jti string) error {
	if s.rdb != nil {
		ctx := context.Background()
		return s.rdb.Del(ctx, jtiKey(jti)).Err()
	}
	// In-memory fallback
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.memStore, jti)
	return nil
}

// RevokeSession revokes a session in PostgreSQL and deletes its JTI from Redis
func (s *TokenStorageService) RevokeSession(ctx context.Context, jti string) error {
	// 1. Remove from Redis
	_ = s.RevokeJTI(jti)

	// 2. Revoke in PostgreSQL
	if s.pgPool != nil {
		query := `UPDATE sessions SET is_revoked = true, updated_at = NOW() WHERE jti = $1`
		_, err := s.pgPool.Exec(ctx, query, jti)
		if err != nil {
			return fmt.Errorf("failed to revoke session in db: %w", err)
		}
	}
	return nil
}

// RevokeSessionBySessionID finds a session by its ID, revokes it in DB, and removes its JTI from Redis
func (s *TokenStorageService) RevokeSessionBySessionID(ctx context.Context, sessionID string) error {
	if s.pgPool != nil {
		var jti string
		err := s.pgPool.QueryRow(ctx, `SELECT jti FROM sessions WHERE session_id = $1`, sessionID).Scan(&jti)
		if err == nil {
			_ = s.RevokeJTI(jti) // Remove refresh token from Redis
		}

		query := `UPDATE sessions SET is_revoked = true, updated_at = NOW() WHERE session_id = $1`
		_, err = s.pgPool.Exec(ctx, query, sessionID)
		if err != nil {
			return fmt.Errorf("failed to revoke session by session id in db: %w", err)
		}
	}
	return nil
}
