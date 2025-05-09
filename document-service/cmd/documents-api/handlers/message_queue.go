package handlers

import (
	"document-service/internal/domain/interfaces"
	"fmt"
)

type MessageHandler struct {
	service interfaces.DocumentServiceInterface
}

func NewMessageHandler(service interfaces.DocumentServiceInterface) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

type MessageHandlerInterface interface {
	HandleMessage(message []byte) error
}

func (m *MessageHandler) HandleMessage(message []byte) error {
	fmt.Println("Received message:", string(message))
	return nil
}
