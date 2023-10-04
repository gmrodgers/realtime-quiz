package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DirectQueue struct {
	Host  string
	Queue string
}

func (dq *DirectQueue) Send(ctx context.Context, msg []byte) error {
	ch, conn, err := createQueue(dq.Host, dq.Queue)
	if err != nil {
		return fmt.Errorf("Failed to declare Queue '%v'", err)
	}
	defer conn.Close()
	defer ch.Close()

	err = ch.PublishWithContext(ctx,
		"",       // exchange
		dq.Queue, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})

	if err != nil {
		return fmt.Errorf("Failed to publish a message")
	}

	return nil
}

func (dq *DirectQueue) Receive(ctx context.Context, until func([][]byte) bool) ([][]byte, error) {
	ch, conn, err := createQueue(dq.Host, dq.Queue)
	if err != nil {
		return [][]byte{}, fmt.Errorf("Failed to declare Queue '%v'", err)
	}
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.ConsumeWithContext(ctx,
		dq.Queue, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return [][]byte{}, fmt.Errorf("Failed to register a consumer: '%v'", err)
	}

	var responses [][]byte
	for {
		select {
		case msg := <-msgs:
			responses = append(responses, msg.Body)
			if until(responses) {
				return responses, nil
			}
		case <-ctx.Done():
			return responses, nil
		}
	}
}

func createQueue(host, name string) (*amqp.Channel, *amqp.Connection, error) {
	conn, err := amqp.Dial(host)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open a channel")
	}

	_, err = ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to declare a queue")
	}

	return ch, conn, nil
}
