package integrationtests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
	"github.com/garrettladley/generate_coding_challenge_server_go/storage"
)

func TestApplicant_ReturnsA200ForValidNUIDThatExistsWithCorrectSolution(t *testing.T) {
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

	var submitResponseBody handlers.SubmitResponseBody

	if err := json.NewDecoder(submitResp.Body).Decode(&submitResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if !submitResponseBody.Correct {
		t.Errorf("Expected solution to be correct, but got: %v", submitResponseBody.Correct)
	}

	if submitResponseBody.Message != "Correct - nice work!" {
		t.Errorf("Expected message to be 'Correct - nice work!', but got: %v", submitResponseBody.Message)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, nuid.String()), nil)

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", resp.StatusCode)
	}

	var applicantResponseBody storage.ApplicantResult

	if err := json.NewDecoder(resp.Body).Decode(&applicantResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if applicantResponseBody.NUID != *nuid {
		t.Errorf("Expected NUID to be %s, but got: %s", nuid.String(), applicantResponseBody.NUID)
	}

	if applicantResponseBody.Name != "Garrett" {
		t.Errorf("Expected name to be Garrett, but got: %s", applicantResponseBody.Name)
	}

	if !*applicantResponseBody.Correct {
		t.Errorf("Expected solution to be correct, but got: %v", applicantResponseBody.Correct)
	}
}

func TestApplicant_ReturnsA200ForValidNUIDThatExistsWithIncorrectSolution(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	nuid, err := domain.ParseNUID("002172052")

	if err != nil {
		t.Errorf("Failed to parse NUID: %v", err)
	}

	registerResp, err := RegisterSampleApplicantWithNUID(app, *nuid)

	if registerResp.HttpStatus != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", registerResp.HttpStatus)
	}

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	submitResp, err := SubmitSolution(app, registerResp, []string{})

	if err != nil {
		t.Errorf("Failed to submit solution: %v", err)
	}

	if submitResp.StatusCode != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", submitResp.StatusCode)
	}

	var submitResponseBody handlers.SubmitResponseBody

	if err := json.NewDecoder(submitResp.Body).Decode(&submitResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if submitResponseBody.Correct {
		t.Errorf("Expected solution to be correct, but got: %v", submitResponseBody.Correct)
	}

	if submitResponseBody.Message != "Incorrect Solution" {
		t.Errorf("Expected message to be 'Incorrect Solution', but got: %v", submitResponseBody.Message)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, nuid.String()), nil)

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", resp.StatusCode)
	}

	var applicantResponseBody storage.ApplicantResult

	if err := json.NewDecoder(resp.Body).Decode(&applicantResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if applicantResponseBody.NUID != *nuid {
		t.Errorf("Expected NUID to be %s, but got: %s", nuid.String(), applicantResponseBody.NUID)
	}

	if applicantResponseBody.Name != "Garrett" {
		t.Errorf("Expected name to be Garrett, but got: %s", applicantResponseBody.Name)
	}

	if *applicantResponseBody.Correct {
		t.Errorf("Expected solution to be incorrect, but got: %v", applicantResponseBody.Correct)
	}
}

func TestApplicant_ReturnsA400ForInvalidNUID(t *testing.T) {
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

	if resp.StatusCode != 400 {
		t.Errorf("Expected status code to be 400, but got: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	if responseString != fmt.Sprintf("invalid NUID %s", badNUID) {
		t.Errorf("Expected response body to be 'invalid NUID %s', but got: %s", badNUID, responseString)
	}
}

func TestApplicant_ReturnsA404ForValidNUIDThatDoesNotExistInDB(t *testing.T) {
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

	if resp.StatusCode != 404 {
		t.Errorf("Expected status code to be 404, but got: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	if responseString != fmt.Sprintf("Applicant with NUID %s not found!", nonexistentNUID) {
		t.Errorf("Expected response body to be 'Applicant with NUID %s not found!', but got: %s", nonexistentNUID, responseString)
	}
}
