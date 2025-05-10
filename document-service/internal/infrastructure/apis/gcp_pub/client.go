package message_queue

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

// MessageHandlerInterface This is a dependency of what is outside the library.
type MessageHandlerInterface interface {
	HandleDocumentCreated(message []byte) error //Example Handler
	HandleAnotherEvent(message []byte) error    // Another handler
}

// PubSubPublisher defines the structure for the Pub/Sub publisher.
type PubSubPublisher struct {
	client    *pubsub.Client
	topicID   string
	projectID string
}

// NewPubSubPublisher creates a new PubSubPublisher.
func NewPubSubPublisher(pubSubClient *pubsub.Client, projectID string, topicID string) *PubSubPublisher {
	return &PubSubPublisher{
		projectID: projectID,
		topicID:   topicID,
		client:    pubSubClient,
	}
}

// Publish publishes a message to Pub/Sub.
func (p *PubSubPublisher) Publish(ctx context.Context, message []byte) (string, error) {
	t := p.client.Topic(p.topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: message,
	})
	msgID, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("Error publishing message: %w", err)
	}
	log.Printf("Published message with ID: %s\n", msgID)
	return msgID, nil
}

// Close closes the Pub/Sub connection.
func (p *PubSubPublisher) Close() {
	if p.client != nil {
		err := p.client.Close()
		if err != nil {
			log.Printf("Error closing Pub/Sub client: %s", err)
		} else {
			log.Println("Pub/Sub client closed.")
		}
	} else {
		log.Println("No Pub/Sub client to close.")
	}
}
