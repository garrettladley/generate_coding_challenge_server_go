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
	"github.com/stretchr/testify/assert"
)

func TestSubmit_ReturnsA200ForCorrectSolution(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	submitResp, err := SubmitCorrectSolution(app)

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
}

func TestSubmit_ReturnsA200ForIncorrectSolution(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	registerResp, err := RegisterSampleApplicant(app)

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
}

func TestSubmit_ReturnsA400ForInvalidToken(t *testing.T) {
	assert := assert.New(t)
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

	assert.Equal(400, submitResp.StatusCode)

	body, err = io.ReadAll(submitResp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	assert.Equal(fmt.Sprintf("invalid token %s", invalidToken), responseString)
}

func TestSubmit_ReturnsIncorrectAfterSubmittingCorrectSolutionThenIncorrectSolution(t *testing.T) {
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

	var submitCorrectResponseBody handlers.SubmitResponseBody

	if err := json.NewDecoder(submitResp.Body).Decode(&submitCorrectResponseBody); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.True(submitCorrectResponseBody.Correct)

	assert.Equal("Correct - nice work!", submitCorrectResponseBody.Message)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/forgot_token/%s", app.Address, nuid), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	forgotTokenResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	assert.Equal(200, forgotTokenResp.StatusCode)

	token, err := GetTokenFromResponse(forgotTokenResp)

	if err != nil {
		t.Errorf("Failed to get token from response: %v", err)
	}

	submitResp, err = SubmitSolutionWithToken(app, *token, []string{})

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
}
