package tests

import (
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/stretchr/testify/assert"
)

func TestGenerateChallenge(t *testing.T) {
	assert := assert.New(t)

	red, _ := domain.Red.String()
	orange, _ := domain.Orange.String()
	yellow, _ := domain.Yellow.String()
	green, _ := domain.Green.String()
	blue, _ := domain.Blue.String()
	violet, _ := domain.Violet.String()

	mandatoryCases := []string{
		"",
		red,
		orange,
		yellow,
		green,
		blue,
		violet,
	}
	nMandatory := len(mandatoryCases)
	nRandom := 10
	challenge := domain.GenerateChallenge(nRandom, mandatoryCases)

	assert.Equal(nMandatory+nRandom, len(challenge.Challenge))

	for _, soln := range challenge.Solution {
		result, err := domain.OneEditAway(soln)
		assert.Nil(err)
		assert.NotNil(result)
	}
}

func TestOneEditAwayExample(t *testing.T) {
	assert := assert.New(t)

	result, err := domain.OneEditAway("red")

	assert.Nil(err)
	assert.Equal(domain.Red, *result)

	result, err = domain.OneEditAway("lue")

	assert.Nil(err)
	assert.Equal(domain.Blue, *result)

	result, err = domain.OneEditAway("ooran")

	assert.EqualErrorf(err, "no valid answer found", "error message %s", "formatted")
	assert.Nil(result)

	result, err = domain.OneEditAway("ooran")

	assert.EqualErrorf(err, "no valid answer found", "error message %s", "formatted")
	assert.Nil(result)

	result, err = domain.OneEditAway("greene")

	assert.Nil(err)
	assert.Equal(domain.Green, *result)
}
