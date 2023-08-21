package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
)

func TestSubmit_ReturnsA200ForCorrectSolution(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	submitResp, err := SubmitCorrectSolution(app)

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
}

func TestSubmit_ReturnsA200ForIncorrectSolution(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	registerResp, err := RegisterSampleApplicant(app)

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
		t.Errorf("Expected solution to be incorrect, but got: %v", submitResponseBody.Correct)
	}

	if submitResponseBody.Message != "Incorrect Solution" {
		t.Errorf("Expected message to be 'Incorrect Solution', but got: %v", submitResponseBody.Message)
	}
}

func TestSubmit_ReturnsA400ForInvalidToken(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	invalidToken := "foo"

	body, err := json.Marshal([]string{})

	if err != nil {
		t.Errorf("Failed to marshal body: %v", err)
	}

	req := httptest.NewRequest("POST", fmt.Sprintf("%s/submit/%s", app.Address, invalidToken), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	submitResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	if submitResp.StatusCode != 400 {
		t.Errorf("Expected status code to be 400, but got: %v", submitResp.StatusCode)
	}

	body, err = io.ReadAll(submitResp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	if responseString != fmt.Sprintf("invalid token %s", invalidToken) {
		t.Errorf("Expected response body to be 'invalid token %s', but got: %v", invalidToken, responseString)
	}
}

func TestSubmit_ReturnsIncorrectAfterSubmittingCorrectSolutionThenIncorrectSolution(t *testing.T) {
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

	var submitCorrectResponseBody handlers.SubmitResponseBody

	if err := json.NewDecoder(submitResp.Body).Decode(&submitCorrectResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if !submitCorrectResponseBody.Correct {
		t.Errorf("Expected solution to be correct, but got: %v", submitCorrectResponseBody.Correct)
	}

	if submitCorrectResponseBody.Message != "Correct - nice work!" {
		t.Errorf("Expected message to be 'Correct - nice work!', but got: %v", submitCorrectResponseBody.Message)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/forgot_token/%s", app.Address, nuid), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	forgotTokenResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	if forgotTokenResp.StatusCode != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", forgotTokenResp.StatusCode)
	}

	token, err := GetTokenFromResponse(forgotTokenResp)

	if err != nil {
		t.Errorf("Failed to get token from response: %v", err)
	}

	submitResp, err = SubmitSolutionWithToken(app, *token, []string{})

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
		t.Errorf("Expected solution to be incorrect, but got: %v", submitResponseBody.Correct)
	}

	if submitResponseBody.Message != "Incorrect Solution" {
		t.Errorf("Expected message to be 'Incorrect Solution', but got: %v", submitResponseBody.Message)
	}
}
