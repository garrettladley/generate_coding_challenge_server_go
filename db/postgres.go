package db

import (
	"github.com/garrettladley/generate_coding_challenge_server_go/config"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

func CreatePostgresConnection(lc fx.Lifecycle, settings config.Settings) *sqlx.DB {
	db := sqlx.MustConnect("postgres", settings.Database.WithDb())

	err := db.Ping()

	if err != nil {
		panic(err)
	}

	return db
}
