package word

import (
	"context"
	"time"
)

type WordConfirmMessage struct {
	UserId   string `bson:"user_id" json:"user_id"`
	ObjectId string `bson:"object_id" json:"object_id"`
}

type WordConfirmedSuccessMessage struct {
	ObjectId    string    `bson:"object_id" json:"object_id"`
	ConfirmedAt time.Time `bson:"confirmed_at" json:"confirmed_at"`
}

type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error

	Close() error
}
