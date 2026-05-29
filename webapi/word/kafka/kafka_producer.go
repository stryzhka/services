package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	metric *prometheus.HistogramVec
}

func NewProducer(brokers []string, topic string, metric *prometheus.HistogramVec) *Producer {

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}
	return &Producer{writer: writer, metric: metric}
}

func (p *Producer) Publish(ctx context.Context, topic, key string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	timer := prometheus.NewTimer(p.metric.WithLabelValues(topic))
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: data,
	})
	log.Println(fmt.Sprintf("publishing %s", key))
	timer.ObserveDuration()
	return err
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
