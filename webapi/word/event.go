package word

import "context"

type WordConfirmMessage struct {
	UserId   string `bson:"user_id" json:"user_id"`
	ObjectId string `bson:"object_id" json:"object_id"`
}

type WordConfirmedSuccessMessage struct {
}

type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error

	Close() error
}
