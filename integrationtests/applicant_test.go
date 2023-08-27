package integrationtests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
)

func TestApplicant_ReturnsA200ForValidNUIDThatExistsWithCorrectSolution(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	nuid, err := domain.ParseNUID("002172052")

	if err != nil {
		t.Errorf("Failed to parse NUID: %v", err)
	}

	submitResp, err := SubmitCorrectSolutionWithNUID(app, *nuid)

	if err != nil {
		t.Errorf("Failed to submit solution: %v", err)
	}

	assert.Equal(200, submitResp.StatusCode)

	var submitResponseBody handlers.SubmitResponseBody

	if err := json.NewDecoder(submitResp.Body).Decode(&submitResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.True(submitResponseBody.Correct)

	assert.Equal("Correct - nice work!", submitResponseBody.Message)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, nuid.String()), nil)

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}

	assert.Equal(200, resp.StatusCode)

	var applicantResponseBody handlers.ApplicantResponse

	if err := json.NewDecoder(resp.Body).Decode(&applicantResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.Equal(*nuid, applicantResponseBody.NUID)

	assert.Equal("Garrett", applicantResponseBody.ApplicantName.String())

	assert.True(applicantResponseBody.Correct)
}

func TestApplicant_ReturnsA200ForValidNUIDThatExistsWithIncorrectSolution(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	nuid, err := domain.ParseNUID("002172052")

	if err != nil {
		t.Errorf("Failed to parse NUID: %v", err)
	}

	registerResp, err := RegisterSampleApplicantWithNUID(app, *nuid)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	submitResp, err := SubmitSolution(app, registerResp, []string{})

	if err != nil {
		t.Errorf("Failed to submit solution: %v", err)
	}

	assert.Equal(200, submitResp.StatusCode)

	var submitResponseBody handlers.SubmitResponseBody

	if err := json.NewDecoder(submitResp.Body).Decode(&submitResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.False(submitResponseBody.Correct)

	assert.Equal("Incorrect Solution", submitResponseBody.Message)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, nuid.String()), nil)

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}

	assert.Equal(200, resp.StatusCode)

	var applicantResponseBody handlers.ApplicantResponse

	if err := json.NewDecoder(resp.Body).Decode(&applicantResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.Equal(*nuid, applicantResponseBody.NUID)

	assert.Equal("Garrett", applicantResponseBody.ApplicantName.String())

	assert.False(applicantResponseBody.Correct)
}

func TestApplicant_ReturnsA400ForInvalidNUID(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	badNUID := "foo"

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, badNUID), nil)

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}

	assert.Equal(400, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	assert.Equal(fmt.Sprintf("invalid NUID %s", badNUID), responseString)
}

func TestApplicant_ReturnsA404ForValidNUIDThatDoesNotExistInDB(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	nonexistentNUID, err := domain.ParseNUID("002172052")

	if err != nil {
		t.Errorf("Failed to parse NUID: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, nonexistentNUID.String()), nil)

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}

	assert.Equal(404, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	assert.Equal(fmt.Sprintf("Applicant with NUID %s not found!", nonexistentNUID.String()), responseString)
}
