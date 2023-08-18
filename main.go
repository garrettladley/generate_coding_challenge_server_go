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
	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func newFiberServer(lc fx.Lifecycle, applicantHandlers *handlers.ApplicantHandler, adminHandlers *handlers.AdminHandler) *fiber.App {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	userGroup := app.Group("/applicant")
	userGroup.Post("/register", applicantHandlers.Register)
	userGroup.Get("/forgot_token/:nuid", applicantHandlers.ForgotToken)
	userGroup.Get("/challenge/:token", applicantHandlers.Challenge)
	userGroup.Post("/submit", applicantHandlers.Submit)

	adminGroup := app.Group("/admin")
	adminGroup.Get("/applicants", adminHandlers.Applicants)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			fmt.Println("Starting fiber server on port 8080")
			go app.Listen(":8080")
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
			config.LoadEnv,
			storage.NewAdminStorage,
			handlers.NewAdminHandler,
			db.CreateMySqlConnection,
			storage.NewApplicantStorage,
			handlers.NewApplicantHandler,
			db.CreateRedisConnection,
		),
		fx.Invoke(newFiberServer),
	).Run()
}
