package message_queue

import (
	"document-service/cmd/documents-api/handlers"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	rabbitMQURL string
	handler     handlers.MessageHandlerInterface
}

func NewRabbitMQConsumer(url string, handler handlers.MessageHandlerInterface) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		rabbitMQURL: url,
		handler:     handler,
	}
}

func (r *RabbitMQConsumer) Consume(queueName string) error {
	// Declare a queue
	q, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Error al declarar la cola: %s", err)
		log.Printf("Cola '%s' declarada. Configurando consumidor...", q.Name)
	}

	err = r.channel.Qos(
		1,     // prefetchCount: número de mensajes a pre-buscar
		0,     // prefetchSize: tamaño en bytes (0 = ilimitado)
		false, // global: aplicar a todos los consumidores en el canal (true) o solo a este (false)
	)

	if err != nil {
		return fmt.Errorf("Error al configurar QoS: %s", err)
	}

	log.Printf("QoS configurado: prefetchCount = %d", 1)

	// 5. Empezar a consumir mensajes
	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer (auto-generado)
		false,  // autoAck: false para confirmación manual
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("Error al consumir mensajes: %s", err)
	}
	log.Printf("Consumidor iniciado para la cola '%s'. Esperando mensajes...", q.Name)

	go func() {
		log.Println("Goroutine de procesamiento de mensajes iniciada", len(msgs))
		for d := range msgs {
			log.Printf("Mensaje recibido")

			handler := r.GetHandlerByConsumer(queueName)
			errHandleMessage := handler(d.Body)
			if errHandleMessage != nil {
				log.Printf("Error al procesar el mensaje: %s", errHandleMessage)
				err = d.Nack(false, false)
				if err != nil {
					log.Printf("Error al no confirmar mensaje: %s", errHandleMessage)
				} else {
					log.Printf("Mensaje no confirmado y reenviado a la cola.")
				}
				continue
			} else {
				log.Printf("Mensaje procesado correctamente.")
			}

			err := d.Ack(false)
			if err != nil {
				log.Printf("Error al confirmar mensaje: %s", err)
			} else {
				log.Printf("Mensaje confirmado.")
			}
		}
		log.Println("Goroutine de procesamiento de mensajes detenida.")
	}()
	select {}
}

func (r *RabbitMQConsumer) Stop() {
	if r.connection != nil {
		err := r.connection.Close()
		if err != nil {
			log.Printf("Error al cerrar la conexión de RabbitMQ: %s", err)
		} else {
			log.Println("Conexión de RabbitMQ cerrada.")
		}
	} else {
		log.Println("No hay conexión de RabbitMQ para cerrar.")
	}
}

func (r *RabbitMQConsumer) Connect() error {
	var conn *amqp.Connection
	var err error
	maxRetries := 5
	retryInterval := 30 * time.Second
	for i := 0; i < maxRetries; i++ {
		time.Sleep(retryInterval)
		log.Printf("Intentando conectar a RabbitMQ (intento %d/%d) en %s", i+1, maxRetries, r.rabbitMQURL)
		conn, err = amqp.Dial(r.rabbitMQURL)
		if err == nil {
			log.Println("Conexión a RabbitMQ establecida.")
			r.connection = conn
			r.channel, err = conn.Channel()
			return nil
		}
		log.Printf("Fallo al conectar a RabbitMQ: %s. Reintentando en %s...", err, retryInterval)
		time.Sleep(retryInterval)
	}
	return fmt.Errorf("Fallo al conectar a RabbitMQ después de %d intentos: %w", maxRetries, err)
}
func (r *RabbitMQConsumer) GetHandlerByConsumer(consumerName string) func(message []byte) error {
	switch consumerName {
	case "register_documents_queue":
		return r.handler.HandleDocumentsRegister
	case "delete_documents_queue":
		return r.handler.HandleDeleteDirectory

	}
	return nil
}
