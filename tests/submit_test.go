package tests

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

	assert.Nil(err)

	submitResp, err := SubmitCorrectSolution(app)

	assert.Nil(err)

	assert.Equal(200, submitResp.StatusCode)

	var submitResponseBody handlers.SubmitResponseBody

	err = json.NewDecoder(submitResp.Body).Decode(&submitResponseBody)

	assert.Nil(err)

	assert.True(submitResponseBody.Correct)

	assert.Equal("Correct - nice work!", submitResponseBody.Message)
}

func TestSubmit_ReturnsA200ForIncorrectSolution(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	registerResp, err := RegisterSampleApplicant(app)

	assert.Nil(err)

	submitResp, err := SubmitSolution(app, registerResp, []string{})

	assert.Nil(err)

	assert.Equal(200, submitResp.StatusCode)

	var submitResponseBody handlers.SubmitResponseBody

	err = json.NewDecoder(submitResp.Body).Decode(&submitResponseBody)

	assert.Nil(err)

	assert.False(submitResponseBody.Correct)

	assert.Equal("Incorrect Solution", submitResponseBody.Message)
}

func TestSubmit_ReturnsA400ForInvalidToken(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	invalidToken := "foo"

	body, err := json.Marshal([]string{})

	assert.Nil(err)

	req := httptest.NewRequest("POST", fmt.Sprintf("%s/submit/%s", app.Address, invalidToken), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	assert.Nil(err)

	submitResp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(400, submitResp.StatusCode)

	body, err = io.ReadAll(submitResp.Body)

	assert.Nil(err)

	responseString := string(body)

	assert.Equal(fmt.Sprintf("invalid token %s", invalidToken), responseString)
}

func TestSubmit_ReturnsIncorrectAfterSubmittingCorrectSolutionThenIncorrectSolution(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	nuid, err := domain.ParseNUID("002172052")

	assert.Nil(err)

	submitResp, err := SubmitCorrectSolutionWithNUID(app, *nuid)

	assert.Nil(err)

	assert.Equal(200, submitResp.StatusCode)

	var submitCorrectResponseBody handlers.SubmitResponseBody

	err = json.NewDecoder(submitResp.Body).Decode(&submitCorrectResponseBody)

	assert.Nil(err)

	assert.True(submitCorrectResponseBody.Correct)

	assert.Equal("Correct - nice work!", submitCorrectResponseBody.Message)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/forgot_token/%s", app.Address, nuid), nil)

	assert.Nil(err)

	forgotTokenResp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(200, forgotTokenResp.StatusCode)

	token, err := GetTokenFromResponse(forgotTokenResp)

	if err != nil {
		t.Errorf("Failed to get token from response: %v", err)
	}

	submitResp, err = SubmitSolutionWithToken(app, *token, []string{})

	assert.Nil(err)

	assert.Equal(200, submitResp.StatusCode)

	var submitResponseBody handlers.SubmitResponseBody

	err = json.NewDecoder(submitResp.Body).Decode(&submitResponseBody)

	assert.Nil(err)

	assert.False(submitResponseBody.Correct)

	assert.Equal("Incorrect Solution", submitResponseBody.Message)
}
