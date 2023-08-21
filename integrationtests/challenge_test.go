package integrationtests

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestChallenge_ReturnsA200ForTokenThatExists(t *testing.T) {
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

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, registerResp.Token), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	challengeResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	if challengeResp.StatusCode != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", challengeResp.StatusCode)
	}

	challengeRespChallenge, err := GetChallengeFromResponse(challengeResp)

	if err != nil {
		t.Errorf("Failed to get challenge from response: %v", err)
	}

	if len(registerResp.Challenge) != len(challengeRespChallenge) {
		t.Errorf("Expected challenge to be %v, but got: %v", registerResp.Challenge, challengeRespChallenge)
	}

	equal := true
	for i := range registerResp.Challenge {
		if registerResp.Challenge[i] != challengeRespChallenge[i] {
			equal = false
			break
		}
	}

	if !equal {
		t.Errorf("Expected challenge to be %v, but got: %v", registerResp.Challenge, challengeRespChallenge)
	}
}

func TestChallenge_ReturnsA400ForInvalidToken(t *testing.T) {
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

	invalidToken := "foo"

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, invalidToken), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	challengeResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	if challengeResp.StatusCode != 400 {
		t.Errorf("Expected status code to be 200, but got: %v", challengeResp.StatusCode)
	}

	body, err := io.ReadAll(challengeResp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	if responseString != fmt.Sprintf("invalid token %s", invalidToken) {
		t.Errorf("Expected response body to be 'invalid token %s', but got: %v", invalidToken, responseString)
	}
}

func TestChallenge_ReturnsA404ForTokenThatDoesNotExistInDB(t *testing.T) {
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, uuid.New()), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	challengeResp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	if challengeResp.StatusCode != 404 {
		t.Errorf("Expected status code to be 200, but got: %v", challengeResp.StatusCode)
	}

	body, err := io.ReadAll(challengeResp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	if responseString != "\"Record associated with given token not found!\"" {
		t.Errorf("Expected response body to be 'Record associated with given token not found!', but got: %v", responseString)
	}
}
