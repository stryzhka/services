package word

import "context"

type WordConfirmMessage struct {
	UserId string `bson:"user_id" json:"user_id"`
}

type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
	Close() error
}
