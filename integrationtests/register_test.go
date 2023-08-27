package integrationtests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type RegisterDB struct {
	ApplicantName sql.NullString      `db:"applicant_name"`
	NUID          sql.NullString      `db:"nuid"`
	Token         sql.NullString      `db:"token"`
	Challenge     storage.StringArray `db:"challenge"`
}

func TestRegister_ReturnsA200ForValidRequestBody(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	resp, err := RegisterSampleApplicant(app)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	numRandom := 100
	numMandatory := 7

	assert.Equal(numRandom+numMandatory, len(resp.Challenge))

	var dbResult RegisterDB

	err = app.Conn.Get(&dbResult, "SELECT applicant_name, nuid, token, challenge FROM applicants;")

	if err != nil {
		t.Errorf("Failed to query database: %v", err)
	}

	assert.True(dbResult.ApplicantName.Valid)
	assert.True(dbResult.NUID.Valid)
	assert.True(dbResult.Token.Valid)
	assert.True(len(dbResult.Challenge) > 0)

	name, err := domain.ParseApplicantName(dbResult.ApplicantName.String)

	if err != nil {
		t.Errorf("Failed to parse applicant name due to invalid database state: %v", err)
	}

	nuid, err := domain.ParseNUID(dbResult.NUID.String)

	if err != nil {
		t.Errorf("Failed to parse NUID due to invalid database state: %v", err)
	}

	token, err := uuid.Parse(dbResult.Token.String)

	if err != nil {
		t.Errorf("Failed to parse token due to invalid database state: %v", err)
	}

	assert.Equal("Garrett", name.String())

	assert.Equal("002172052", nuid.String())

	assert.Equal(token, resp.Token)

	assert.Equal(len(resp.Challenge), len(dbResult.Challenge))

	for i := 0; i < len(resp.Challenge); i++ {
		assert.Equal(resp.Challenge[i], dbResult.Challenge[i])
	}
}

func TestRegister_ReturnsA400WhenRequestBodyPropertiesAreMissing(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	testCases := make([]map[string]string, 3)

	testCases[0] = map[string]string{
		"name": "Garrett",
	}

	testCases[1] = map[string]string{
		"nuid": "002172052",
	}

	testCases[2] = map[string]string{}

	for _, testCase := range testCases {
		body, err := json.Marshal(testCase)

		if err != nil {
			t.Errorf("Failed to marshal request body: %v", err)
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("%s/register", app.Address), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.App.Test(req)

		if err != nil {
			t.Errorf("Failed to register applicant: %v", err)
		}

		assert.Equal(400, resp.StatusCode)

		var dbResult RegisterDB

		err = app.Conn.Get(&dbResult, "SELECT applicant_name, nuid FROM applicants;")

		assert.True(err != nil)

		assert.False(dbResult.ApplicantName.Valid)

		assert.False(dbResult.NUID.Valid)

		assert.True(err == sql.ErrNoRows)
	}
}

func TestRegister_ReturnsA400WhenRequestBodyPropertiesArePresentButInvalid(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	testCases := make([]map[string]string, 3)

	testCases[0] = map[string]string{
		"name": "",
		"nuid": "002172052",
	}

	testCases[1] = map[string]string{
		"name": "Garrett",
		"nuid": "",
	}

	testCases[2] = map[string]string{"name": "", "nuid": ""}

	for _, testCase := range testCases {
		body, err := json.Marshal(testCase)

		if err != nil {
			t.Errorf("Failed to marshal request body: %v", err)
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("%s/register", app.Address), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.App.Test(req)

		if err != nil {
			t.Errorf("Failed to register applicant: %v", err)
		}

		assert.Equal(400, resp.StatusCode)

		var dbResult RegisterDB

		err = app.Conn.Get(&dbResult, "SELECT applicant_name, nuid FROM applicants;")

		assert.True(err != nil)

		assert.False(dbResult.ApplicantName.Valid)

		assert.False(dbResult.NUID.Valid)

		assert.True(err == sql.ErrNoRows)
	}
}

func TestRegister_ReturnsA409ForUserThatAlreadyExists(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	nuid, err := domain.ParseNUID("002172052")

	if err != nil {
		t.Errorf("Failed to parse NUID: %v", err)
	}

	_, err = RegisterSampleApplicantWithNUID(app, *nuid)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	resp, err := RegisterRequest(app, *nuid)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	assert.Equal(409, resp.StatusCode)

	var count int

	err = app.Conn.Get(&count, "SELECT COUNT(*) FROM applicants WHERE nuid = $1;", nuid.String())

	if err != nil {
		t.Errorf("Failed to query database: %v", err)
	}

	assert.Equal(1, count)
}
