package main

import (
	"context"

	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/config"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/services"
)

var (
	ctx context.Context
	err error

	agendaSvc services.AgendaService
)

func init() {
	ctx = context.TODO()
	config.LoadConfig()
}

func main() {
	ctx = context.Background()

	// Starts the agenda service
	agendaSvc, err = services.NewAgendaService()
	if err != nil {
		panic(err)
	}
	go agendaSvc.StartBroker()

	go agendaSvc.AgendaTask()
	// Waits for the context to be done
	<-ctx.Done()
}
