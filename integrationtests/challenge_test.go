package integrationtests

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChallenge_ReturnsA200ForTokenThatExists(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	registerResp, err := RegisterSampleApplicant(app)

	if err != nil {
		t.Errorf("Failed to register applicant: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, registerResp.Token), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	challengeResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	assert.Equal(200, challengeResp.StatusCode)

	challengeRespChallenge, err := GetChallengeFromResponse(challengeResp)

	if err != nil {
		t.Errorf("Failed to get challenge from response: %v", err)
	}

	assert.Equal(len(registerResp.Challenge), len(challengeRespChallenge))

	for i := 0; i < len(registerResp.Challenge); i++ {
		assert.Equal(registerResp.Challenge[i], challengeRespChallenge[i])
	}
}

func TestChallenge_ReturnsA400ForInvalidToken(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	invalidToken := "foo"

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, invalidToken), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	challengeResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	assert.Equal(400, challengeResp.StatusCode)

	body, err := io.ReadAll(challengeResp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	assert.Equal(fmt.Sprintf("invalid token %s", invalidToken), responseString)
}

func TestChallenge_ReturnsA404ForTokenThatDoesNotExistInDB(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	nonexistentToken := uuid.New().String()

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, nonexistentToken), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	challengeResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	assert.Equal(404, challengeResp.StatusCode)

	body, err := io.ReadAll(challengeResp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	assert.Equal(fmt.Sprintf("Record associated with token %s not found!", nonexistentToken), responseString)
}
