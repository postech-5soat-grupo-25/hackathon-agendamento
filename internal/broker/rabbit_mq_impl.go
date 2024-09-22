package broker

import (
	"context"
	"errors"

	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
)

type RabbitMqImpl struct {
}

func (r *RabbitMqImpl) Consume(ctx context.Context, consumerChan chan *models.Appointment) error {
	return errors.New("off")
}

func NewBroker() (Broker, error) {
	return &RabbitMqImpl{}, nil
}
