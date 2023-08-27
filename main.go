package main

import (
	"github.com/garrettladley/generate_coding_challenge_server_go/config"
	"github.com/garrettladley/generate_coding_challenge_server_go/db"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
	"github.com/garrettladley/generate_coding_challenge_server_go/server"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"

	"go.uber.org/fx"
)

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
		fx.Invoke(server.NewFxFiberApp),
	).Run()
}
