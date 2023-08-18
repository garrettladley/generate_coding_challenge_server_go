package db

import (
	"github.com/garrettladley/generate_coding_challenge_server_go/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func CreatePostgresConnection(lc fx.Lifecycle, config config.Settings) *sqlx.DB {
	db := sqlx.MustConnect("postgres", config.Database.WithDb())

	err := db.Ping()
	if err != nil {
		panic(err)
	} else {
		println("DB CONNECTED")
	}

	return db
}
