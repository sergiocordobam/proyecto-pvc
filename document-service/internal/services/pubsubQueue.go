package services

import (
	"context"
	"document-service/internal/domain/interfaces"
	"document-service/internal/domain/models"
	"encoding/json"

	"github.com/labstack/gommon/log"
)

type SendNotificationsService struct {
	client interfaces.PubSubClientInterface
}

func NewSendNotificationsService(client interfaces.PubSubClientInterface) *SendNotificationsService {
	return &SendNotificationsService{
		client: client,
	}
}

func (p *SendNotificationsService) SendNotification(ctx context.Context, notificationRequest models.NotificationMessage) error {
	messageBytes, marshalMessageErr := json.Marshal(notificationRequest)
	if marshalMessageErr != nil {
		return marshalMessageErr
	}
	msg, err := p.client.Publish(ctx, messageBytes)
	if err != nil {
		return err
	}
	log.Info(ctx, "Message of action published successfully: ", msg)
	return nil
}
