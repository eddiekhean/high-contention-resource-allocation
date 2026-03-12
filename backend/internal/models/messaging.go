package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ─── PostgreSQL Entities ─────────────────────────────────────────────────────

type PGUser struct {
	UserID       string    `db:"user_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	DisplayName  string    `db:"display_name"`
	AvatarURL    string    `db:"avatar_url"`
	Timezone     string    `db:"timezone"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Device struct {
	DeviceID    string     `db:"device_id"`
	UserID      string     `db:"user_id"`
	PushToken   string     `db:"push_token"`
	DeviceModel string     `db:"device_model"`
	LastActive  *time.Time `db:"last_active"`
	IsRevoked   bool       `db:"is_revoked"`
	CreatedAt   time.Time  `db:"created_at"`
}

type ConnectionStatus string

const (
	ConnectionPending  ConnectionStatus = "pending"
	ConnectionAccepted ConnectionStatus = "accepted"
	ConnectionRejected ConnectionStatus = "rejected"
	ConnectionBlocked  ConnectionStatus = "blocked"
)

type Connection struct {
	ConnectionID string           `db:"connection_id"`
	RequesterID  string           `db:"requester_id"`
	AddresseeID  string           `db:"addressee_id"`
	Status       ConnectionStatus `db:"status"`
	CreatedAt    time.Time        `db:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at"`
}

type ConversationType string

const (
	ConversationPair  ConversationType = "pair"
	ConversationGroup ConversationType = "group"
)

type Conversation struct {
	ConversationID string           `db:"conversation_id"`
	Type           ConversationType `db:"type"`
	CreatedAt      time.Time        `db:"created_at"`
}

type ConversationMember struct {
	ConversationID string     `db:"conversation_id"`
	UserID         string     `db:"user_id"`
	Role           string     `db:"role"` // "member" or "admin"
	JoinedAt       time.Time  `db:"joined_at"`
	LeftAt         *time.Time `db:"left_at"`
}

type E2EEKeyType string

const (
	KeyTypeIdentity     E2EEKeyType = "identity"
	KeyTypePrekey       E2EEKeyType = "prekey"
	KeyTypeOnetimePrekey E2EEKeyType = "onetime_prekey"
)

type E2EEKey struct {
	KeyID         string      `db:"key_id"`
	UserID        string      `db:"user_id"`
	DeviceID      string      `db:"device_id"`
	PublicKeyData string      `db:"public_key_data"`
	KeyType       E2EEKeyType `db:"key_type"`
	IsActive      bool        `db:"is_active"`
	CreatedAt     time.Time   `db:"created_at"`
}

type Session struct {
	SessionID  string    `db:"session_id"`
	UserID     string    `db:"user_id"`
	DeviceID   *string   `db:"device_id"`
	JTI        string    `db:"jti"`
	IPAddress  *string   `db:"ip_address"`
	UserAgent  *string   `db:"user_agent"`
	ExpiresAt  time.Time `db:"expires_at"`
	IsRevoked  bool      `db:"is_revoked"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// ─── MongoDB Documents ────────────────────────────────────────────────────────

type DeliveryMode string

const (
	DeliveryRealtime  DeliveryMode = "REALTIME"
	DeliveryTimeShift DeliveryMode = "TIMESHIFT"
)

type MessageStatus string

const (
	MessageScheduled MessageStatus = "scheduled"
	MessageReleased  MessageStatus = "released"
	MessageDelivered MessageStatus = "delivered"
	MessageRead      MessageStatus = "read"
)

type MessageContent struct {
	Ciphertext string `bson:"ciphertext"`
	IV         string `bson:"iv"`
	KeyID      string `bson:"keyId"`
}

type MessageTimestamps struct {
	CreatedAtUtc          time.Time  `bson:"createdAtUtc"`
	RecipientVisibleAtUtc time.Time  `bson:"recipientVisibleAtUtc"`
	ReleasedAtUtc         *time.Time `bson:"releasedAtUtc"`
	DeliveredAt           *time.Time `bson:"deliveredAt"`
	ReadAt                *time.Time `bson:"readAt"`
}

type MessageMetadata struct {
	RecipientTimeZone  string `bson:"recipientTimeZone"`
	RecipientLocalTime string `bson:"recipientLocalTime"`
}

type Message struct {
	ID             bson.ObjectID     `bson:"_id,omitempty"`
	ConversationID string            `bson:"conversationId"`
	SenderID       string            `bson:"senderId"`
	MessageType    string            `bson:"messageType"` // "text", "image", etc.
	DeliveryMode   DeliveryMode      `bson:"deliveryMode"`
	Content        MessageContent    `bson:"content"`
	Timestamps     MessageTimestamps `bson:"timestamps"`
	Status         MessageStatus     `bson:"status"`
	Metadata       MessageMetadata   `bson:"metadata"`
}
