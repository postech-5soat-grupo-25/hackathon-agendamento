package services

import (
	"context"
	"fmt"
	"log/slog"
	"encoding/json"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/broker"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/storage"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/controllers"
)

type AgendaServiceImpl struct {
	queue chan *models.Message

	ctx    context.Context
	broker broker.Broker
	db     storage.Storage
	controller controllers.AppointmentsInterface
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
			// API
			//slog.Log(a.ctx, slog.LevelInfo, "Appointment scheduled")
			fmt.Println(msg.Type)
			switch msg.Type {
				case models.GetDoctorWorkingHoursMessageType:
					var message models.GetDoctorWorkingHoursMessage
					err := json.Unmarshal(msg.Body, &message)
					if err != nil {
						return fmt.Errorf("failed to unmarshal GetDoctorWorkingHoursMessageType message: %v", err)
					}
					response := a.controller.GetWorkingHours(message.DoctorID)

					fmt.Println(response.StatusCode)
					fmt.Println(string(response.Body))

				case models.WorkingHoursMessage:
					var message models.WorkingHours
					err := json.Unmarshal(msg.Body, &message)
					if err != nil {
						return fmt.Errorf("failed to unmarshal WorkingHours message: %v", err)
					}
					response := a.controller.CreateOrEditWorkingHours(&message)

					fmt.Println(response.StatusCode)
					fmt.Println(string(response.Body))

				case models.GetClientAppointmentsMessage:
					var message models.GenericIDMessage
					err := json.Unmarshal(msg.Body, &message)
					if err != nil {
						return fmt.Errorf("failed to unmarshal GetClientAppointmentsMessage message: %v", err)
					}
					response := a.controller.GetClientAppointments(message.ID)
					fmt.Println(response.StatusCode)
					fmt.Println(string(response.Body))

				case models.AppointmentMessage:
					var message models.Appointment
					err := json.Unmarshal(msg.Body, &message)
					if err != nil {
						return fmt.Errorf("failed to unmarshal AppointmentMessage message: %v", err)
					}
					response := a.controller.ScheduleAppointment(&message)
					fmt.Println(response.StatusCode)
					fmt.Println(string(response.Body))

				case models.CancelScheduledAppointmentMessage:
					var message models.GenericIDMessage
					err := json.Unmarshal(msg.Body, &message)
					if err != nil {
						return fmt.Errorf("failed to unmarshal CancelScheduledAppointmentMessage message: %v", err)
					}
					response := a.controller.CancelScheduledAppointment(message.ID)
					fmt.Println(response.StatusCode)
					fmt.Println(string(response.Body))
				
				default:
					fmt.Println("Unrecognized message type:" + msg.Type)
			}
	
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

	controller := controllers.NewController(db)
	
	slog.Log(ctx, slog.LevelDebug, "storage connected")
	fmt.Println("Done with Connections")
	return &AgendaServiceImpl{
		make(chan *models.Message),
		ctx,
		broker,
		db,
		controller,
	}, nil
}
