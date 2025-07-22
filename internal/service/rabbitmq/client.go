package rabbitmq

import (
	"adel/internal/service/postgres"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QUEUE_NAME = "submission_queue"

type RabbitMQClient struct {
	conn   *amqp.Connection
	logger *log.Logger
}

type SubmissionJob struct {
	Submission postgres.Submission `json:"submission"`
	Problem    postgres.Problem    `json:"problem"`
}

func Open(logger *log.Logger) *RabbitMQClient {
	conn, err := amqp.Dial("amqp://admin:adminpassword@localhost:5672/")
	if err != nil {
		logger.Panicf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.Panicf("failed to open a channel for queue declaration: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		QUEUE_NAME, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		logger.Panicf("failed to declare queue: %v", err)
	}

	return &RabbitMQClient{
		conn:   conn,
		logger: logger,
	}
}

func (rc *RabbitMQClient) Close() {
	err := rc.conn.Close()
	if err != nil {
		rc.logger.Printf("failed to close connection: %v", err)
	}
}

func (rc *RabbitMQClient) PublishSubmissionJob(job *SubmissionJob) error {
	ch, err := rc.conn.Channel()
	if err != nil {
		rc.logger.Printf("failed to open channel: %v", err)
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(job)
	if err != nil {
		rc.logger.Printf("failed to marshal job: %v", err)
		return err
	}

	err = ch.Publish(
		"",         // exchange
		QUEUE_NAME, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		rc.logger.Printf("failed to publish message: %v", err)
		return err
	}

	return nil
}

func (rc *RabbitMQClient) ConsumeSubmissionJobs() (*amqp.Channel, <-chan amqp.Delivery, error) {
	ch, err := rc.conn.Channel()
	if err != nil {
		rc.logger.Printf("failed to open channel: %v", err)
		return nil, nil, err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		rc.logger.Printf("failed to set QoS: %v", err)
		return nil, nil, err
	}

	msgs, err := ch.Consume(
		QUEUE_NAME, // queue
		"",         // consumer tags
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		ch.Close()
		rc.logger.Printf("failed to consume messages: %v", err)
		return nil, nil, err
	}

	return ch, msgs, nil
}
