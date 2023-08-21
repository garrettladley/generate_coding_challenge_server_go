package domain

import (
	"testing"
)

func TestGenerateChallenge(t *testing.T) {
	red, _ := Red.String()
	orange, _ := Orange.String()
	yellow, _ := Yellow.String()
	green, _ := Green.String()
	blue, _ := Blue.String()
	violet, _ := Violet.String()

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
	challenge := GenerateChallenge(nRandom, mandatoryCases)

	if len(challenge.Challenge) != nMandatory+nRandom {
		t.Errorf("Expected challenge length: %d, got: %d", nMandatory+nRandom, len(challenge.Challenge))
	}

	for _, soln := range challenge.Solution {
		result, err := OneEditAway(soln)
		resultStr, _ := result.String()
		if resultStr == "" || err != nil {
			t.Errorf("Expected valid answer, but got %v %v", resultStr, err)
		}
	}
}

func TestOneEditAwayExample(t *testing.T) {

	if result, err := OneEditAway("red"); result != Red || err != nil {
		t.Errorf("Expected Color red, but got %v %v", result, err)
	}
	if result, err := OneEditAway("lue"); result != Blue || err != nil {
		t.Errorf("Expected Color blue, but got %v %v", result, err)
	}
	if result, err := OneEditAway("ooran"); err == nil {
		t.Errorf("Expected nil, but got %v %v", result, err)
	}
	if result, err := OneEditAway("abc"); err == nil {
		t.Errorf("Expected nil, but got %v %v", result, err)
	}
	if result, err := OneEditAway("greene"); result != Green || err != nil {
		t.Errorf("Expected Color Green, but got %v %v", result, err)
	}
}
