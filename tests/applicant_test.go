package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/garrettladley/generate_coding_challenge_server_go/handlers"
)

func TestApplicant_ReturnsA200ForValidNUIDThatExistsWithCorrectSolution(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	nuid, err := domain.ParseNUID("002172052")

	assert.Nil(err)

	submitResp, err := SubmitCorrectSolutionWithNUID(app, *nuid)

	assert.Nil(err)
	assert.Equal(200, submitResp.StatusCode)

	var submitResponseBody handlers.SubmitResponseBody

	err = json.NewDecoder(submitResp.Body).Decode(&submitResponseBody)

	assert.Nil(err)

	assert.True(submitResponseBody.Correct)

	assert.Equal("Correct - nice work!", submitResponseBody.Message)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, nuid.String()), nil)

	resp, err := app.App.Test(req)

	assert.Nil(err)
	assert.Equal(200, resp.StatusCode)

	var applicantResponseBody handlers.ApplicantResponse

	err = json.NewDecoder(resp.Body).Decode(&applicantResponseBody)

	assert.Nil(err)

	assert.Equal(*nuid, applicantResponseBody.NUID)

	assert.Equal("Garrett", applicantResponseBody.ApplicantName.String())

	assert.True(applicantResponseBody.Correct)
}

func TestApplicant_ReturnsA200ForValidNUIDThatExistsWithIncorrectSolution(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	nuid, err := domain.ParseNUID("002172052")

	assert.Nil(err)

	registerResp, err := RegisterSampleApplicantWithNUID(app, *nuid)

	assert.Nil(err)

	submitResp, err := SubmitSolution(app, registerResp, []string{})

	assert.Nil(err)
	assert.Equal(200, submitResp.StatusCode)

	var submitResponseBody handlers.SubmitResponseBody

	err = json.NewDecoder(submitResp.Body).Decode(&submitResponseBody)

	assert.Nil(err)

	assert.False(submitResponseBody.Correct)

	assert.Equal("Incorrect Solution", submitResponseBody.Message)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, nuid.String()), nil)

	resp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(200, resp.StatusCode)

	var applicantResponseBody handlers.ApplicantResponse

	err = json.NewDecoder(resp.Body).Decode(&applicantResponseBody)

	assert.Nil(err)

	assert.Equal(*nuid, applicantResponseBody.NUID)

	assert.Equal("Garrett", applicantResponseBody.ApplicantName.String())

	assert.False(applicantResponseBody.Correct)
}

func TestApplicant_ReturnsA400ForInvalidNUID(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	badNUID := "foo"

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, badNUID), nil)

	resp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(400, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)

	assert.Nil(err)

	responseString := string(body)

	assert.Equal(fmt.Sprintf("invalid NUID %s", badNUID), responseString)
}

func TestApplicant_ReturnsA404ForValidNUIDThatDoesNotExistInDB(t *testing.T) {
	assert := assert.New(t)
	app, err := SpawnApp()

	assert.Nil(err)

	nonexistentNUID, err := domain.ParseNUID("002172052")

	assert.Nil(err)

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/applicant/%s", app.Address, nonexistentNUID.String()), nil)

	resp, err := app.App.Test(req)

	assert.Nil(err)

	assert.Equal(404, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)

	assert.Nil(err)

	responseString := string(body)

	assert.Equal(fmt.Sprintf("Applicant with NUID %s not found!", nonexistentNUID.String()), responseString)
}
