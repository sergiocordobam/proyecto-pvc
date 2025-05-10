package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// GenericMessage represents the generic JSON message structure from other APIs.
type GenericMessage struct {
	Event     string                 `json:"event"`
	User      int                    `json:"user"`
	Name      string                 `json:"name"`
	UserEmail string                 `json:"user_email"`
	ExtraData map[string]interface{} `json:"extra_data"`
}

// PubSubConfig defines the configuration for Google Cloud Pub/Sub.
type PubSubConfig struct {
	ProjectID       string
	SubscriptionIDs []string // Changed to a list of subscriptions
}

// SendEmail sends an email using SendGrid
func SendEmail(fromEmail string, apiKey string, message GenericMessage) error {
	fmt.Sprintf("Sending email to %s\n", message.UserEmail)
	from := mail.NewEmail("PVC Documentos", fromEmail)
	subject := fmt.Sprintf("Notification: %s", message.Event)
	to := mail.NewEmail(message.Name, message.UserEmail)

	// Build the email body dynamically
	var body strings.Builder
	body.WriteString(fmt.Sprintf("Event: %s\n", message.Event))
	body.WriteString(fmt.Sprintf("User ID: %d\n", message.User))

	// Iterate over extra data
	for key, value := range message.ExtraData {
		body.WriteString(fmt.Sprintf("%s: %v\n", key, formatValue(value)))
	}
	fmt.Println("Body: ", body.String())

	plainTextContent := body.String()
	htmlContent := fmt.Sprintf("<strong>%s</strong>", strings.ReplaceAll(body.String(), "\n", "<br>")) //example HTML
	messageSG := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(messageSG)
	if err != nil {
		return fmt.Errorf("sendgrid send error: %w", err)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid API error: %d, %s", response.StatusCode, response.Body)
	}
	fmt.Sprintf("Email sent to %s with status code %d\n", message.UserEmail, response.StatusCode)

	return nil
}

// formatValue formats the generic value passed to the email sender, converting if needed
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%.2f", v) // Format float with 2 decimal places
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v) // Default string representation
	}
}

// EmailWorker is a worker that listens for email notifications from Pub/Sub and sends them.
func EmailWorker(ctx context.Context, pubsubClient *pubsub.Client, pubSubConfig PubSubConfig, fromEmail string, sendGridApiKey string, subscriptionID string) {
	sub := pubsubClient.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("Received message ID: %s from subscription: %s\n", msg.ID, subscriptionID)

		// Unmarshal the message data into a GenericMessage
		var genericMessage GenericMessage
		err := json.Unmarshal(msg.Data, &genericMessage)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			msg.Nack() // Indicate that the message processing failed so that will be retried.  Don't Ack!
			return     // Exit the message handler
		}
		fmt.Sprintf("Message: %s\n", genericMessage)

		err = SendEmail(fromEmail, sendGridApiKey, genericMessage)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			msg.Nack()
			return
		}
		log.Println("Message: ", genericMessage)

		msg.Ack() // Acknowledge the message only after successful processing
	})

	if err != nil {
		log.Printf("Error receiving messages from subscription %s: %v", subscriptionID, err)
	}
}

func main() {
	ctx := context.Background()

	// Configuration - Replace with your actual values
	pubSubConfig := PubSubConfig{
		ProjectID: "zeta-matrix-458323-p1", // Replace
		SubscriptionIDs: []string{
			"auth-api-topic-sub",      //Replace
			"documents-topic-sub",     //Replace
			"interoperator-topic-sub", //Replace
		},
	}

	fromEmail := os.Getenv("FROM_EMAIL")            // Set your From Email address in env variable
	sendGridApiKey := os.Getenv("SENDGRID_API_KEY") // Set your SendGrid API key in env variable
	fmt.Println("FROM_EMAIL: ", fromEmail)
	if fromEmail == "" || sendGridApiKey == "" {
		log.Fatalf("FROM_EMAIL and SENDGRID_API_KEY environment variables must be set")
	}

	// Create a pubsub client
	pubsubClient, err := pubsub.NewClient(ctx, pubSubConfig.ProjectID)
	if err != nil {
		log.Fatalf("Could not create pubsub client: %v", err)
	}
	defer pubsubClient.Close()

	// Start multiple email workers (subscribers), one for each subscription
	var wg sync.WaitGroup
	wg.Add(len(pubSubConfig.SubscriptionIDs)) // the wait group counts the subscripctions.
	for _, subID := range pubSubConfig.SubscriptionIDs {
		go func(subscriptionID string) { // Capture the subscriptionID in the goroutine
			defer wg.Done()
			EmailWorker(ctx, pubsubClient, pubSubConfig, fromEmail, sendGridApiKey, subscriptionID)
		}(subID) // Pass the subscriptionID to the goroutine
	}
	wg.Wait()
	fmt.Println("All email workers finished")
}
