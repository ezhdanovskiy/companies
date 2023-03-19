package kafka

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type AsyncProducer struct {
	writer *kafka.Writer
}

type ProducerConfig struct {
	Brokers      []string
	Topic        string
	BatchSize    int
	BatchTimeout time.Duration
}

func NewAsyncProducer(cfg *ProducerConfig) *AsyncProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
		RequiredAcks: kafka.RequireNone,
		Async:        true,
		Balancer:     &kafka.Murmur2Balancer{},
		Logger:       kafka.LoggerFunc(zap.S().Debugf),
		ErrorLogger:  kafka.LoggerFunc(zap.S().Errorf),
	}
	return &AsyncProducer{
		writer: writer,
	}
}

func (ap *AsyncProducer) Publish(ctx context.Context, messages ...[]byte) error {
	mm := make([]kafka.Message, 0, len(messages))

	for i := range messages {
		message := messages[i]
		mm = append(mm, kafka.Message{
			Value: message,
		})
	}

	if len(mm) > 0 {
		err := ap.writer.WriteMessages(ctx, mm...)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (ap *AsyncProducer) Close() error {
	return ap.writer.Close()
}
