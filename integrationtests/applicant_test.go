package integrationtests

import (
	"encoding/json"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
)

func TestApplicant_ReturnsA200ForValidNUIDThatExistsWithCorrectSolution(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	registerResp, err := RegisterSampleApplicant(app)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	if registerResp.HttpStatus != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", registerResp.HttpStatus)
	}

	solution := []string{}

	for _, challenge := range registerResp.Challenge {
		result, err := domain.OneEditAway(challenge)
		if err == nil {
			result, err := result.String()
			if err != nil {
				t.Errorf("Failed to convert result to string: %v", err)
			}
			solution = append(solution, result)
		}
	}

	submitResp, err := SubmitSolution(registerResp, app, solution)

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

	if !submitResponseBody.Correct {
		t.Errorf("Expected solution to be correct, but got: %v", submitResponseBody.Correct)
	}

	if submitResponseBody.Message != "Correct - nice work!" {
		t.Errorf("Expected message to be 'Correct - nice work!', but got: %v", submitResponseBody.Message)
	}
}

func TestApplicant_ReturnsA200ForValidNUIDThatExistsWithIncorrectSolution(t *testing.T) {
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

	submitResp, err := SubmitSolution(registerResp, app, []string{})

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
}
