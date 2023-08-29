package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/stretchr/testify/assert"
)

func TestParseApplicantName_ValidName(t *testing.T) {
	assert := assert.New(t)

	nameInput := "Muneer Lalji"
	result, err := domain.ParseApplicantName(nameInput)

	assert.Nil(err)
	assert.Equal(domain.ApplicantName(nameInput), &result)
}

func TestParseApplicantName_256GraphemeLongNameIsValid(t *testing.T) {
	assert := assert.New(t)

	nameInput := strings.Repeat("a", 256)
	result, err := domain.ParseApplicantName(nameInput)

	assert.Nil(err)
	assert.Equal(domain.ApplicantName(nameInput), &result)
}

func TestParseApplicantName_NameLongerThan256GraphemesIsRejected(t *testing.T) {
	assert := assert.New(t)

	nameInput := strings.Repeat("a", 257)
	result, err := domain.ParseApplicantName(nameInput)

	assert.Errorf(err, fmt.Sprintf("invalid name! Given: %s", nameInput), "error message %s", "formatted")
	assert.Nil(result)
}

func TestParseApplicantName_WhitespaceOnlyNamesAreRejected(t *testing.T) {
	assert := assert.New(t)

	nameInput := " "
	result, err := domain.ParseApplicantName(nameInput)

	assert.Errorf(err, fmt.Sprintf("invalid name! Given: %s", nameInput), "error message %s", "formatted")
	assert.Nil(result)
}

func TestParseApplicantName_EmptyStringIsRejected(t *testing.T) {
	assert := assert.New(t)

	nameInput := ""
	result, err := domain.ParseApplicantName(nameInput)

	assert.Errorf(err, fmt.Sprintf("invalid name! Given: %s", nameInput), "error message %s", "formatted")
	assert.Nil(result)
}

func TestParseApplicantName_NamesContainingInvalidCharactersAreRejected(t *testing.T) {
	assert := assert.New(t)

	invalidCharacters := []string{"/", "(", ")", "\"", "<", ">", "\\", "{", "}"}

	for _, char := range invalidCharacters {
		nameInput := char
		result, err := domain.ParseApplicantName(nameInput)

		assert.Errorf(err, fmt.Sprintf("invalid name! Given: %s", nameInput), "error message %s", "formatted")
		assert.Nil(result)
	}
}
