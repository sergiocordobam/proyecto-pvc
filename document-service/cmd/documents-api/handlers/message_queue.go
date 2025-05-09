package handlers

import (
	"document-service/internal/domain/interfaces"
	"document-service/internal/domain/models"
	"encoding/json"
	"errors"
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
	var requestRegisterDocs models.RegisterDocumentsMessage
	if err := json.Unmarshal(message, &requestRegisterDocs); err != nil {
		return err
	}
	if requestRegisterDocs.CitizenId == 0 {
		return errors.New("CitizenId is required")
	}
	if len(requestRegisterDocs.Documents) == 0 {
		return errors.New("Documents are required")
	}
	return nil
}
