package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string) *Consumer {
	return &Consumer{kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     topic,
			MaxBytes:    10e6, // 10MB
			StartOffset: kafka.FirstOffset,
		})}
}

func (c *Consumer) Consume(ctx context.Context, handler func(key, value []byte) error) error {
	log.Println("consumer started, waiting for messages...") // ← добавь
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("read error: %v", err) // ← и это
			return err
		}
		log.Printf("got message: %s", string(msg.Value)) // ← и это
		if err := handler(msg.Key, msg.Value); err != nil {
			return err
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
