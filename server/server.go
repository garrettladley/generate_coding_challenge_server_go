package server

import (
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func NewFiberApp(address string, applicantHandlers *handlers.ApplicantHandler, adminHandlers *handlers.AdminHandler) *fiber.App {
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

	go app.Listen(address)

	return app
}
