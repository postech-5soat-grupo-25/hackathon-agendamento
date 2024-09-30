package main

import (
	"context"
	"log/slog"
	"github.com/gofiber/fiber/v2"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/models"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/storage"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/controllers"
	"github.com/postech-5soat-grupo-25/hackathon-agendamento/internal/config"
)

var (
	ctx context.Context
	err error
	app    fiber.App
	db     storage.Storage
	controller controllers.AppointmentsInterface
)

func init() {
	ctx = context.TODO()
	config.LoadConfig()
}

func main() {
	ctx = context.Background()

	db, err := storage.NewStorage(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		return
	}

	controller := controllers.NewController(db)
	
	slog.Log(ctx, slog.LevelDebug, "storage connected")

	app := fiber.New()

    // Health check endpoint with log
    app.Get("/health", func(c *fiber.Ctx) error {
        slog.Log(ctx, slog.LevelInfo, "Health check endpoint hit")
        return c.SendStatus(200)
    })

    app.Get("/working-hours", func(c *fiber.Ctx) error {
        var msg models.GetDoctorWorkingHoursMessage
        if err := c.QueryParser(&msg); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid query parameters"})
        }
        response := controller.GetWorkingHours(msg.DoctorID)
        return c.Status(response.StatusCode).JSON(response)
    })

    app.Post("/working-hours", func(c *fiber.Ctx) error {
        var workhours models.WorkingHours
        if err := c.BodyParser(&workhours); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }
        response := controller.CreateOrEditWorkingHours(&workhours)
        return c.Status(response.StatusCode).JSON(response)
    })

    app.Get("/appointment", func(c *fiber.Ctx) error {
        var msg models.GenericIDMessage
        if err := c.QueryParser(&msg); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid query parameters"})
        }
        response := controller.GetClientAppointments(msg.ID)
        return c.Status(response.StatusCode).JSON(response)
    })

    app.Post("/appointment", func(c *fiber.Ctx) error {
        var appointment models.Appointment
        if err := c.BodyParser(&appointment); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }
        response := controller.ScheduleAppointment(&appointment)
        return c.Status(response.StatusCode).JSON(response)
    })

    app.Delete("/appointment", func(c *fiber.Ctx) error {
        var msg models.GenericIDMessage
        if err := c.BodyParser(&msg); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
        }
        response := controller.CancelScheduledAppointment(msg.ID)
        return c.Status(response.StatusCode).JSON(response)
    })

	app.Listen(":80")
}
