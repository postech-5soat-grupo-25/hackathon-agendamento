package broker

import (
	"context"

	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
)

type Broker interface {
	Consume(ctx context.Context, consumerChan chan *models.Message) error
	Close() error
}
