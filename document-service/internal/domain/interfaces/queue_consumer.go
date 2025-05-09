package interfaces

type QueueConsumerInterface interface {
	Consume(queueName string) error
	Connect() error
	Stop()
}
