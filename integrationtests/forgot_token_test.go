package integrationtests

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
)

func TestForgot_Token_ReturnsA200ForNUIDThatExists(t *testing.T) {
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

	if registerResp.HttpStatus != 200 {
		t.Errorf("Expected status code to be 200, but got: %v", registerResp.HttpStatus)
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

	if token.String() != registerResp.Token.String() {
		t.Errorf("Expected token to be %s, but got: %v", registerResp.Token, token)
	}
}

func TestForgot_Token_ReturnssA400ForInvalidNUID(t *testing.T) {
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

func TestForgot_Token_ReturnssA400ForNUIDThatDoesNotExistInDB(t *testing.T) {
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
