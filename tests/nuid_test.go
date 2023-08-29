package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/stretchr/testify/assert"
)

func TestParseNUID_ValidNUID(t *testing.T) {
	assert := assert.New(t)

	nuidInput := strings.Repeat("1", 9)
	result, err := domain.ParseNUID(nuidInput)

	assert.Nil(err)
	assert.Equal(domain.NUID(nuidInput), *result)
}

func TestParseNUID_WhitespaceOnlyIsRejected(t *testing.T) {
	assert := assert.New(t)

	nuidInput := " "
	result, err := domain.ParseNUID(nuidInput)

	assert.Errorf(err, fmt.Sprintf("invalid NUID! Given: %s", nuidInput), "error message %s", "formatted")
	assert.Nil(result)
}

func TestParseNUID_EmptyStringIsRejected(t *testing.T) {
	assert := assert.New(t)

	nuidInput := ""
	result, err := domain.ParseNUID(nuidInput)

	assert.Errorf(err, fmt.Sprintf("invalid NUID! Given: %s", nuidInput), "error message %s", "formatted")
	assert.Nil(result)
}

func TestParseNUID_InvalidLengthIsRejected(t *testing.T) {
	assert := assert.New(t)

	nuidInput := strings.Repeat("1", 10)
	result, err := domain.ParseNUID(nuidInput)

	assert.Errorf(err, fmt.Sprintf("invalid NUID! Given: %s", nuidInput), "error message %s", "formatted")
	assert.Nil(result)
}

func TestParseNUID_InvalidCharactersAreRejected(t *testing.T) {
	assert := assert.New(t)

	nuidInput := "a" + strings.Repeat("1", 8)
	result, err := domain.ParseNUID(nuidInput)

	assert.Errorf(err, fmt.Sprintf("invalid NUID! Given: %s", nuidInput), "error message %s", "formatted")
	assert.Nil(result)
}

func TestParseNUID_StringWithInvalidCharactersIsRejected(t *testing.T) {
	assert := assert.New(t)

	characters := []rune{'1', 'a'}

	for numA := 1; numA <= 8; numA++ {
		permutation := strings.Repeat("a", numA)
		fullString := permutation + "11111111"[0:8-numA]

		for i := 0; i < 9; i++ {
			for _, char := range characters {
				testString := domain.InsertCharAtIndex(fullString, char, i)
				result, err := domain.ParseNUID(testString)

				assert.Errorf(err, fmt.Sprintf("invalid NUID! Given: %s", testString), "error message %s", "formatted")
				assert.Nil(result)
			}
		}
	}
}
