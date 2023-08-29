package tests

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

	assert.Nil(err)

	registerResp, err := RegisterSampleApplicant(app)

	assert.Nil(err)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, registerResp.Token), nil)

	assert.Nil(err)

	challengeResp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(200, challengeResp.StatusCode)

	challengeRespChallenge, err := GetChallengeFromResponse(challengeResp)

	assert.Nil(err)

	assert.Equal(len(registerResp.Challenge), len(challengeRespChallenge))

	for i := 0; i < len(registerResp.Challenge); i++ {
		assert.Equal(registerResp.Challenge[i], challengeRespChallenge[i])
	}
}

func TestChallenge_ReturnsA400ForInvalidToken(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	invalidToken := "foo"

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, invalidToken), nil)

	assert.Nil(err)

	challengeResp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(400, challengeResp.StatusCode)

	body, err := io.ReadAll(challengeResp.Body)

	assert.Nil(err)

	responseString := string(body)

	assert.Equal(fmt.Sprintf("invalid token %s", invalidToken), responseString)
}

func TestChallenge_ReturnsA404ForTokenThatDoesNotExistInDB(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	nonexistentToken := uuid.New().String()

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/challenge/%s", app.Address, nonexistentToken), nil)

	assert.Nil(err)

	challengeResp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(404, challengeResp.StatusCode)

	body, err := io.ReadAll(challengeResp.Body)

	assert.Nil(err)

	responseString := string(body)

	assert.Equal(fmt.Sprintf("Record associated with token %s not found!", nonexistentToken), responseString)
}
