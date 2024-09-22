package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/broker"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
)

type AgendaServiceImpl struct {
	queue chan *models.Appointment

	ctx    context.Context
	broker broker.Broker
}

func (a *AgendaServiceImpl) StartBroker() error {
	go func() {
		if err := a.broker.Consume(a.ctx, a.queue); err != nil {
			slog.Log(a.ctx, slog.LevelError, err.Error())
		}
	}()
	return nil
}

func (a *AgendaServiceImpl) AgendaTask() error {
	for {
		select {
		case msg := <-a.queue:
			// consume message
			slog.Log(a.ctx, slog.LevelInfo, "Appointment scheduled")
			fmt.Println(msg)
			return nil
		case <-a.ctx.Done():
			return a.ctx.Err()
		}
	}
}

func NewAgendaService() (AgendaService, error) {
	broker, err := broker.NewBroker()
	if err != nil {
		return nil, err
	}
	return &AgendaServiceImpl{
		make(chan *models.Appointment),
		context.Background(),
		broker,
	}, nil
}
