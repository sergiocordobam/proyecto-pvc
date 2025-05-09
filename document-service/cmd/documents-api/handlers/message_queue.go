package handlers

import (
	"context"
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
	HandleDocumentsRegister(message []byte) error
	HandleDeleteDirectory(message []byte) error
}

func (m *MessageHandler) HandleDocumentsRegister(message []byte) error {
	ctx := context.Background()
	var requestRegisterDocs models.RegisterDocumentsMessage
	if err := json.Unmarshal(message, &requestRegisterDocs); err != nil {
		return err
	}
	if requestRegisterDocs.CitizenId == 0 {
		return errors.New("CitizenId is required")
	}
	errTransfer := m.service.TransferDocsToCurrentBucket(ctx, requestRegisterDocs)
	if errTransfer != nil {
		return errTransfer
	}
	return nil
}
func (m *MessageHandler) HandleDeleteDirectory(message []byte) error {
	ctx := context.Background()
	var requestDeleteDocuments models.DeleteDocumentsMessage
	if err := json.Unmarshal(message, &requestDeleteDocuments); err != nil {
		return err
	}
	if requestDeleteDocuments.CitizenId == 0 {
		return errors.New("CitizenId is required")
	}
	errDelete := m.service.DeleteAllFilesInUserDirectory(ctx, requestDeleteDocuments.CitizenId)
	if errDelete != nil {
		return errDelete
	}
	return nil
}
