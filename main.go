package main

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"github.com/garrettladley/generate_coding_challenge_server_go/config"
	"github.com/garrettladley/generate_coding_challenge_server_go/db"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
	"github.com/garrettladley/generate_coding_challenge_server_go/server"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"

	"github.com/gofiber/fiber/v2"
)

func newFiberServer(lc fx.Lifecycle, settings config.Settings, applicantHandlers *handlers.ApplicantHandler, adminHandlers *handlers.AdminHandler) *fiber.App {
	address := fmt.Sprintf(":%d", settings.Application.Port)
	app := server.NewFiberApp(address, applicantHandlers, adminHandlers)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go app.Listen(address)
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
