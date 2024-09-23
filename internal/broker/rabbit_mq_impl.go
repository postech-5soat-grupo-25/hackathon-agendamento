package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/streadway/amqp"

	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/config"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
)

var (
	conn *amqp.Connection
	err  error
)

type RabbitMqImpl struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func (r *RabbitMqImpl) Consume(ctx context.Context, consumerChan chan *models.Appointment) error {
	// Set up consumer
	msgs, err := r.channel.Consume(
		r.queue.Name,            // queue
		"hackathon-agendamento", // consumer
		false,                   // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}
	fmt.Println("Waiting for messages. To exit press CTRL+C")
	// Start consuming messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}
				appointment := &models.Appointment{}
				err := json.Unmarshal(d.Body, appointment)
				if err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					d.Nack(false, true)
					continue
				}
				consumerChan <- appointment
				d.Ack(false)
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func (r *RabbitMqImpl) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

func NewBroker(ctx context.Context) (Broker, error) {
	// Get the RabbitMQ URL from environment variable
	amqpURL := config.GetEnvHost()
	if amqpURL == "" {
		return nil, errors.New("RABBITMQ_URL environment variable is not set")
	}

	for {
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			break
		}
		slog.Log(ctx, slog.LevelError, err.Error())
		slog.Log(ctx, slog.LevelError, "Failed to connect to RabbitMQ, retrying in 5 seconds...")
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	// Declare a queue
	q, err := ch.QueueDeclare(
		"appointments", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	return &RabbitMqImpl{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}
