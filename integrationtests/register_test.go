package integrationtests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
)

type RegisterDB struct {
	ApplicantName sql.NullString `db:"applicant_name"`
	Nuid          sql.NullString `db:"nuid"`
}

func TestRegister_ReturnsA200ForValidRequestBody(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	resp, err := RegisterSampleApplicant(app)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	if resp.HttpStatus != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", resp.HttpStatus)
	}

	numRandom := 100
	numMandatory := 7

	if len(resp.Challenge) != numRandom+numMandatory {
		t.Errorf("Expected 'challenge' length to be %v, but got: %v", numRandom+numMandatory, len(resp.Challenge))
	}

	var dbResult RegisterDB

	err = app.Conn.Get(&dbResult, "SELECT applicant_name, nuid FROM applicants;")

	if err != nil {
		t.Errorf("Failed to query database: %v", err)
	}

	if !dbResult.ApplicantName.Valid && !dbResult.Nuid.Valid {
		t.Error("Expected database to contain applicant name and nuid")
	}

	name, err := domain.ParseApplicantName(dbResult.ApplicantName.String)

	if err != nil {
		t.Errorf("Failed to parse applicant name due to invalid database state: %v", err)
	}

	nuid, err := domain.ParseNUID(dbResult.Nuid.String)

	if err != nil {
		t.Errorf("Failed to parse NUID due to invalid database state: %v", err)
	}

	if name.String() != "Garrett" {
		t.Errorf("Expected applicant name to be 'Garrett', but got: %v", name)
	}

	if nuid.String() != "002172052" {
		t.Errorf("Expected NUID to be '002172052', but got: %v", nuid)
	}
}

func TestRegister_ReturnsA400WhenRequestBodyPropertiesAreMissing(t *testing.T) {
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

		if resp.StatusCode != 400 {
			t.Errorf("Expected status code to be 400, but got: %v", resp.StatusCode)
		}

		var dbResult RegisterDB

		err = app.Conn.Get(&dbResult, "SELECT applicant_name, nuid FROM applicants;")

		if err != nil {
			if dbResult.ApplicantName.Valid && dbResult.Nuid.Valid && err != sql.ErrNoRows {
				t.Errorf("Expected database to be empty, but got: %v", err)
			}
		} else if dbResult.ApplicantName.Valid || dbResult.Nuid.Valid {
			t.Errorf("Expected database to be empty, but got: %v", err)
		}
	}
}

func TestRegister_ReturnsA400WhenRequestBodyPropertiesArePresentButInvalid(t *testing.T) {
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

		if resp.StatusCode != 400 {
			t.Errorf("Expected status code to be 400, but got: %v", resp.StatusCode)
		}

		var dbResult RegisterDB

		err = app.Conn.Get(&dbResult, "SELECT applicant_name, nuid FROM applicants;")

		if err != nil {
			if dbResult.ApplicantName.Valid && dbResult.Nuid.Valid && err != sql.ErrNoRows {
				t.Errorf("Expected database to be empty, but got: %v", err)
			}
		} else if dbResult.ApplicantName.Valid || dbResult.Nuid.Valid {
			t.Errorf("Expected database to be empty, but got: %v", err)
		}
	}
}

func TestRegister_ReturnsA409ForUserThatAlreadyExists(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	nuid, err := domain.ParseNUID("002172052")

	if err != nil {
		t.Errorf("Failed to parse NUID: %v", err)
	}

	resp, err := RegisterSampleApplicantWithNUID(app, *nuid)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	if resp.HttpStatus != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", resp.HttpStatus)
	}

	resp, err = RegisterSampleApplicantWithNUID(app, *nuid)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	if resp.HttpStatus != 409 {
		t.Errorf("Expected status code to be 409, but got: %v", resp.HttpStatus)
	}

	var count int

	err = app.Conn.Get(&count, "SELECT COUNT(*) FROM applicants WHERE nuid = $1;", nuid.String())

	if err != nil {
		t.Errorf("Failed to query database: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected database to contain 1 applicant, but got: %v", count)
	}
}
