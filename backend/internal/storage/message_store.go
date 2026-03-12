package storage

import (
	"context"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const messagesCollection = "messages"

// MessageStore handles MongoDB operations for messages.
type MessageStore struct {
	col *mongo.Collection
}

func NewMessageStore(db *mongo.Database) *MessageStore {
	return &MessageStore{col: db.Collection(messagesCollection)}
}

// EnsureIndexes creates the 3 required indexes for the messages collection.
// Safe to call on every startup — mongo skips if index already exists.
func (s *MessageStore) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		// Worker index: find scheduled messages ready to release
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "timestamps.recipientVisibleAtUtc", Value: 1},
			},
			Options: options.Index().SetName("idx_worker_scheduled"),
		},
		// Inbox index: fetch conversation history sorted by visibility time
		{
			Keys: bson.D{
				{Key: "conversationId", Value: 1},
				{Key: "timestamps.recipientVisibleAtUtc", Value: -1},
			},
			Options: options.Index().SetName("idx_inbox_conversation"),
		},
		// Outbox index: sender's sent messages filtered by status
		{
			Keys: bson.D{
				{Key: "senderId", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().SetName("idx_outbox_sender"),
		},
	}

	_, err := s.col.Indexes().CreateMany(ctx, indexes)
	return err
}

// Insert adds a new message document.
func (s *MessageStore) Insert(ctx context.Context, msg *models.Message) error {
	_, err := s.col.InsertOne(ctx, msg)
	return err
}

// FindScheduledReady returns messages whose recipientVisibleAtUtc <= now and status == scheduled.
func (s *MessageStore) FindScheduledReady(ctx context.Context) ([]models.Message, error) {
	filter := bson.M{
		"status": models.MessageScheduled,
		"timestamps.recipientVisibleAtUtc": bson.M{"$lte": time.Now().UTC()},
	}
	cursor, err := s.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var msgs []models.Message
	if err := cursor.All(ctx, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// Release updates a message status to released and sets releasedAtUtc.
func (s *MessageStore) Release(ctx context.Context, id bson.ObjectID) error {
	now := time.Now().UTC()
	_, err := s.col.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"status":                      models.MessageReleased,
			"timestamps.releasedAtUtc":    now,
		},
	})
	return err
}

// FindInbox returns visible messages for a conversation (excludes scheduled TIMESHIFT).
func (s *MessageStore) FindInbox(ctx context.Context, conversationID string, limit int64) ([]models.Message, error) {
	filter := bson.M{
		"conversationId": conversationID,
		"$or": bson.A{
			bson.M{"deliveryMode": models.DeliveryRealtime},
			bson.M{
				"deliveryMode": models.DeliveryTimeShift,
				"status":       bson.M{"$ne": models.MessageScheduled},
			},
		},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamps.recipientVisibleAtUtc", Value: -1}}).
		SetLimit(limit)

	cursor, err := s.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var msgs []models.Message
	if err := cursor.All(ctx, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}
