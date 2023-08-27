package integrationtests

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/stretchr/testify/assert"
)

func TestForgot_Token_ReturnsA200ForNUIDThatExists(t *testing.T) {
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

	assert.Equal(&registerResp.Token, token)
}

func TestForgot_Token_ReturnssA400ForInvalidNUID(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	badNUID := "foo"

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/forgot_token/%s", app.Address, badNUID), nil)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	resp, err := app.App.Test(req)

	if err != nil {
		t.Errorf("Failed to execute request: %v", err)
	}

	assert.Equal(400, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	responseString := string(body)

	assert.Equal(fmt.Sprintf("invalid NUID %s", badNUID), responseString)
}

func TestForgot_Token_ReturnssA400ForNUIDThatDoesNotExistInDB(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	if err != nil {
		t.Errorf("Failed to spawn app: %v", err)
	}

	nonexistentNUID, err := domain.ParseNUID("002172052")

	if err != nil {
		t.Errorf("Failed to parse NUID: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/forgot_token/%s", app.Address, nonexistentNUID), nil)

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

	assert.Equal(fmt.Sprintf("Applicant with NUID %s not found!", nonexistentNUID), responseString)
}
