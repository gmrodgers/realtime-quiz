package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type FanoutQueue struct {
	Host     string
	Exchange string
}

func (fq *FanoutQueue) Broadcast(ctx context.Context, msg []byte) error {
	ch, conn, err := createExchange(fq.Host, fq.Exchange)
	if err != nil {
		return fmt.Errorf("Failed to declare Exchange '%v'", err)
	}
	defer conn.Close()
	defer ch.Close()

	err = ch.PublishWithContext(ctx,
		fq.Exchange, // exchange
		"",          // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})

	if err != nil {
		return fmt.Errorf("Failed to publish a message")
	}

	return nil
}

func (fq *FanoutQueue) Listen(ctx context.Context, until func([][]byte) bool) ([][]byte, error) {
	ch, conn, err := createExchange(fq.Host, fq.Exchange)
	if err != nil {
		return [][]byte{}, fmt.Errorf("Failed to declare Exchange '%v'", err)
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return [][]byte{}, fmt.Errorf("Failed to declare a queue: '%v'", err)
	}

	err = ch.QueueBind(
		q.Name,      // queue name
		"",          // routing key
		fq.Exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return [][]byte{}, fmt.Errorf("Failed to bind queue to exchange: '%v'", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return [][]byte{}, fmt.Errorf("Failed to consume queue: '%v'", err)
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

func createExchange(host, name string) (*amqp.Channel, *amqp.Connection, error) {
	conn, err := amqp.Dial(host)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to RabbitMQ: '%v'", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open a channel: '%v'", err)
	}

	err = ch.ExchangeDeclare(
		name,     // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to declare an exchange: '%v'", err)
	}

	return ch, conn, nil
}
