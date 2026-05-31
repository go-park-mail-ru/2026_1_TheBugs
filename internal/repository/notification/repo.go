package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/pkg/kfk"
	"github.com/segmentio/kafka-go"
)

type KafkaEmailProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaEmailProducer() *KafkaEmailProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   config.NotificationTopic,
	})
	return &KafkaEmailProducer{
		writer: w,
		topic:  config.NotificationTopic,
	}
}
func (p *KafkaEmailProducer) send(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaEmailProducer) SendJSON(ctx context.Context, key string, data interface{}) error {
	value, err := json.Marshal(kfk.Event{Payload: data, Type: string(key)})
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	return p.send(ctx, []byte(key), value)
}

func (p *KafkaEmailProducer) SendAnswer(ctx context.Context, req dto.AnswerNotification) error {
	return p.SendJSON(ctx, entity.AnswerRoute, req)
}
func (p *KafkaEmailProducer) SendRecoveryCode(ctx context.Context, req dto.RecoveryNotification) error {
	return p.SendJSON(ctx, entity.RecoveryRoute, req)
}
func (p *KafkaEmailProducer) SendVerificationCode(ctx context.Context, req dto.VerificationNotification) error {
	return p.SendJSON(ctx, entity.VerificationRoute, req)
}
func (p *KafkaEmailProducer) SendRoommateMatch(ctx context.Context, req dto.RoommateMatchNotification) error {
	return p.SendJSON(ctx, entity.RoommateMatchRoute, req)
}
func (p *KafkaEmailProducer) SendRoommateContactsForRequester(ctx context.Context, req dto.RoommateContactsNotification) error {
	return p.SendJSON(ctx, entity.RoommateContactsForRequesterRoute, req)
}
func (p *KafkaEmailProducer) SendRoommateContactsForAccepted(ctx context.Context, req dto.RoommateContactsNotification) error {
	return p.SendJSON(ctx, entity.RoommateContactsForAcceptedRoute, req)
}

func (p *KafkaEmailProducer) Close() error {
	return p.writer.Close()
}
