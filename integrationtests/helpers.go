package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/garrettladley/generate_coding_challenge_server_go/config"
	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
	"github.com/garrettladley/generate_coding_challenge_server_go/server"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
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
		result[i] = letterBytes[domain.GenerateRandomInt(int64(len(letterBytes)))]
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

func RegisterSampleApplicant(app TestApp) (*storage.RegisterApplicantResponse, error) {
	return RegisterSampleApplicantWithNUID(app, "002172052")
}

func RegisterSampleApplicantWithNUID(app TestApp, nuid domain.NUID) (*storage.RegisterApplicantResponse, error) {
	data := map[string]string{
		"name": "Garrett",
		"nuid": nuid.String(),
	}

	body, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", fmt.Sprintf("%s/register", app.Address), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.App.Test(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return &storage.RegisterApplicantResponse{HttpStatus: resp.StatusCode}, nil
	}

	var responseBody map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	challengeStrings, err := GetChallengeFromBody(responseBody)

	if err != nil {
		return nil, err
	}

	token, err := GetTokenFromBody(responseBody)

	if err != nil {
		return nil, err
	}

	return &storage.RegisterApplicantResponse{
		Challenge:  challengeStrings,
		Token:      token,
		HttpStatus: resp.StatusCode,
	}, nil
}

func GetChallengeFromBody(responseBody map[string]interface{}) ([]string, error) {
	challenge, challengeExists := responseBody["challenge"].([]interface{})

	if !challengeExists {
		return nil, fmt.Errorf("response does not contain 'challenge' property")
	}

	challengeStrings := make([]string, len(challenge))

	for i, v := range challenge {
		challengeStrings[i] = v.(string)
	}

	return challengeStrings, nil
}

func GetTokenFromBody(responseBody map[string]interface{}) (*uuid.UUID, error) {
	token, tokenExists := responseBody["token"]

	if !tokenExists {
		return nil, fmt.Errorf("response does not contain 'token' property")
	}

	parsedToken, err := uuid.Parse(token.(string))

	if err != nil {
		return nil, err
	}

	return &parsedToken, nil
}

func GetChallengeFromResponse(resp *http.Response) ([]string, error) {
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return GetChallengeFromBody(responseBody)
}

func SubmitSolution(registerResponse *storage.RegisterApplicantResponse, app TestApp, solution []string) (*http.Response, error) {
	body, err := json.Marshal(solution)

	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", fmt.Sprintf("%s/submit/%s", app.Address, registerResponse.Token), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.App.Test(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
