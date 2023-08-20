package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"github.com/garrettladley/generate_coding_challenge_server_go/config"
	"github.com/garrettladley/generate_coding_challenge_server_go/db"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func newFiberServer(lc fx.Lifecycle, settings config.Settings, applicantHandlers *handlers.ApplicantHandler, adminHandlers *handlers.AdminHandler) *fiber.App {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}:${port} ${pid} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/register", applicantHandlers.Register)
	app.Get("/forgot_token/:nuid", applicantHandlers.ForgotToken)
	app.Get("/challenge/:token", applicantHandlers.Challenge)
	app.Post("/submit/:token", applicantHandlers.Submit)

	app.Get("/applicant/:nuid", adminHandlers.Applicant)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go app.Listen(fmt.Sprintf(":%d", settings.Application.Port))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}

func main() {
	fx.New(
		fx.Provide(
			config.GetConfiguration,
			db.CreatePostgresConnection,
			storage.NewAdminStorage,
			handlers.NewAdminHandler,
			storage.NewApplicantStorage,
			handlers.NewApplicantHandler,
		),
		fx.Invoke(newFiberServer),
	).Run()
}
