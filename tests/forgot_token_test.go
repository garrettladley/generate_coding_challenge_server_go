package tests

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

	assert.Nil(err)

	nuid, err := domain.ParseNUID("002172052")

	assert.Nil(err)

	registerResp, err := RegisterSampleApplicantWithNUID(app, *nuid)

	assert.Nil(err)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/forgot_token/%s", app.Address, nuid), nil)

	assert.Nil(err)

	forgotTokenResp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(200, forgotTokenResp.StatusCode)

	token, err := GetTokenFromResponse(forgotTokenResp)

	assert.Nil(err)

	assert.Equal(&registerResp.Token, token)
}

func TestForgot_Token_ReturnssA400ForInvalidNUID(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	badNUID := "foo"

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/forgot_token/%s", app.Address, badNUID), nil)

	assert.Nil(err)

	resp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(400, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)

	assert.Nil(err)

	responseString := string(body)

	assert.Equal(fmt.Sprintf("invalid NUID %s", badNUID), responseString)
}

func TestForgot_Token_ReturnssA400ForNUIDThatDoesNotExistInDB(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	nonexistentNUID, err := domain.ParseNUID("002172052")

	assert.Nil(err)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/forgot_token/%s", app.Address, nonexistentNUID), nil)

	resp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(404, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)

	assert.Nil(err)

	responseString := string(body)

	assert.Equal(fmt.Sprintf("Applicant with NUID %s not found!", nonexistentNUID), responseString)
}
