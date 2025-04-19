package queue

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Dsn               string
	ConsumerQueueName string
	ConsumerName      string
	AutoAck           bool
	Args              amqp.Table
	Channel           *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	rabbitMQArgs := amqp.Table{}
	rabbitMQArgs["x-dead-letter-exchange"] = os.Getenv("RABBITMQ_DLX")

	return &RabbitMQ{
		Dsn:               os.Getenv("RABBITMQ_DSN"),
		ConsumerName:      os.Getenv("RABBITMQ_CONSUMER_NAME"),
		ConsumerQueueName: os.Getenv("RABBITMQ_CONSUMER_QUEUE_NAME"),
		AutoAck:           false,
		Args:              amqp.Table{},
	}
}

func (r *RabbitMQ) Connect() *amqp.Channel {
	conn, err := amqp.Dial(r.Dsn)
	failOnError(err, "Failed to connect to RabbitMQ")

	r.Channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	return r.Channel
}

func (r *RabbitMQ) Consume(messageChannel chan amqp.Delivery) {
	q, err := r.Channel.QueueDeclare(
		r.ConsumerQueueName,
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		r.Args, // arguments
	)
	failOnError(err, "Failed to declare a queue")

	incomingMessages, err := r.Channel.Consume(
		q.Name,         // queue
		r.ConsumerName, // consumer
		r.AutoAck,      // auto-acknowledge
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		r.Args,         // arguments
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for message := range incomingMessages {
			log.Println("Incoming new message")
			messageChannel <- message
		}
		log.Println("RabbitMQ channel closed")
		close(messageChannel)
	}()
}

func (r *RabbitMQ) Notify(message string, contentType string, exchange string, routingKey string) error {
	err := r.Channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
