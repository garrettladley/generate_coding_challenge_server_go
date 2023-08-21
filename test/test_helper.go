package test

import (
	"fmt"
	"math/rand"
	"net"
	"os"

	"github.com/garrettladley/generate_coding_challenge_server_go/config"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
	"github.com/garrettladley/generate_coding_challenge_server_go/server"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type TestApp struct {
	App     *fiber.App
	Address string
	Conn    *sqlx.DB
}

func SpawnApp() (TestApp, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		return TestApp{}, err
	}

	initialDir, err := os.Getwd()

	if err != nil {
		return TestApp{}, err
	}

	err = os.Chdir("../")

	if err != nil {
		return TestApp{}, err
	}

	configuration, err := config.GetConfiguration()

	if err != nil {
		return TestApp{}, err
	}

	err = os.Chdir(initialDir)

	if err != nil {
		return TestApp{}, err
	}

	configuration.Database.DatabaseName = generateRandomDBName()

	connectionWithDB, err := configureDatabase(configuration.Database)

	if err != nil {
		fmt.Print("foo")
		return TestApp{}, err
	}

	return TestApp{
		App:     server.NewFiberApp(listener.Addr().String(), handlers.NewApplicantHandler(storage.NewApplicantStorage(connectionWithDB)), handlers.NewAdminHandler(storage.NewAdminStorage(connectionWithDB))),
		Address: fmt.Sprintf("http://%s", listener.Addr().String()),
		Conn:    connectionWithDB,
	}, nil
}

func generateRandomDBName() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyz"
	length := 36
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(result)
}

func configureDatabase(config config.DatabaseSettings) (*sqlx.DB, error) {
	connectionWithoutDB := sqlx.MustConnect("postgres", config.WithoutDb())

	_, err := connectionWithoutDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", config.DatabaseName))

	if err != nil {
		return nil, err
	}

	connectionWithDB := sqlx.MustConnect("postgres", config.WithDb())

	driver, err := postgres.WithInstance(connectionWithDB.DB, &postgres.Config{})

	if err != nil {
		return nil, err
	}

	initialDir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	err = os.Chdir("../")

	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		config.DatabaseName,
		driver,
	)

	if err != nil {
		return nil, err
	}

	err = os.Chdir(initialDir)

	if err != nil {
		return nil, err
	}

	err = m.Up()

	if err != nil {
		return nil, err
	}

	return connectionWithDB, nil
}
