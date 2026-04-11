package kafka

import (
	"context"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	metric *prometheus.HistogramVec
}

func NewConsumer(brokers []string, topic string, metric *prometheus.HistogramVec) *Consumer {
	return &Consumer{kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     topic,
			MaxBytes:    10e6, // 10MB
			StartOffset: kafka.FirstOffset,
		}),
		metric}
}

func (c *Consumer) Consume(ctx context.Context, handler func(key, value []byte) error) error {
	log.Println("consumer started, waiting for messages...")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("read error: %v", err)
			return err
		}
		log.Printf("got message: %s", string(msg.Value))
		timer := prometheus.NewTimer(c.metric.WithLabelValues("consumed"))
		err = handler(msg.Key, msg.Value)
		timer.ObserveDuration()

		if err != nil {
			return err
		}
		if err := handler(msg.Key, msg.Value); err != nil {
			return err
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
