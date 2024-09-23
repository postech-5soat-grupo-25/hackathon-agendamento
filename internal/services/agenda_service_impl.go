package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/broker"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/storage"
)

type AgendaServiceImpl struct {
	queue chan *models.Appointment

	ctx    context.Context
	broker broker.Broker
	db     storage.Storage
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
	ctx := context.Background()
	broker, err := broker.NewBroker(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		return nil, err
	}
	slog.Log(ctx, slog.LevelDebug, "broker connected")

	db, err := storage.NewStorage(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		return nil, err
	}
	slog.Log(ctx, slog.LevelDebug, "storage connected")
	fmt.Println("Done with Connections")
	return &AgendaServiceImpl{
		make(chan *models.Appointment),
		ctx,
		broker,
		db,
	}, nil
}
