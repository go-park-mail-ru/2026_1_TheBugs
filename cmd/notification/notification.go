package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/smtp"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/pkg/kfk"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	err := config.Read(logrus.StandardLogger())
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:29092"},
		Topic:       config.NotificationTopic,
		GroupID:     config.NotificationGroupID,
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()
	router := kfk.NewRouter()

	senderRepo := smtp.NewSMTPSender(
		config.Config.SMTP.Host,
		config.Config.SMTP.Port,
		config.Config.SMTP.Email,
		config.Config.SMTP.Pwd,
	)

	router.Register(entity.AnswerRoute, func(ctx context.Context, payload any) error {
		var data dto.AnswerNotification
		if err := mapstructure.Decode(payload, &data); err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}
		senderRepo.SendAnswer(ctx, data)
		return nil
	})

	router.Register(entity.RecoveryRoute, func(ctx context.Context, payload any) error {
		var data dto.RecoveryNotification
		if err := mapstructure.Decode(payload, &data); err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}
		senderRepo.SendRecoveryCode(ctx, data)
		return nil
	})

	router.Register(entity.VerificationRoute, func(ctx context.Context, payload any) error {
		var data dto.VerificationNotification
		if err := mapstructure.Decode(payload, &data); err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}
		senderRepo.SendVerificationCode(ctx, data)
		return nil
	})

	router.Register(entity.RoommateMatchRoute, func(ctx context.Context, payload any) error {
		var data dto.RoommateMatchNotification
		if err := mapstructure.Decode(payload, &data); err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}
		senderRepo.SendRoommateMatch(ctx, data)
		return nil
	})

	router.Register(entity.RoommateContactsForRequesterRoute, func(ctx context.Context, payload any) error {
		var data dto.RoommateContactsNotification
		if err := mapstructure.Decode(payload, &data); err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}
		senderRepo.SendRoommateContactsForRequester(ctx, data)
		return nil
	})

	router.Register(entity.RoommateContactsForAcceptedRoute, func(ctx context.Context, payload any) error {
		var data dto.RoommateContactsNotification
		if err := mapstructure.Decode(payload, &data); err != nil {
			return fmt.Errorf("failed to decode payload: %w", err)
		}
		senderRepo.SendRoommateContactsForAccepted(ctx, data)
		return nil
	})
	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			continue
		}

		log.Printf("Чтение offset=%d", msg.Offset)
		log.Printf("Получено сообщение: %s", string(msg.Value))
		err = router.Handle(context.Background(), msg.Value)
		if err != nil {
			log.Printf(" router.Handle: %v", err)
			continue
		}
		err = reader.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Printf(" router.Handle: %v", err)
			continue
		}
		log.Printf("Offset %d закоммичен", msg.Offset)
	}
}
