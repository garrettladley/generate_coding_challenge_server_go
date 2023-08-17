package domain

import (
	"strings"
	"testing"
)

func TestParseNUID_ValidNUID(t *testing.T) {
	nuidInput := strings.Repeat("1", 9)
	_, err := ParseNUID(nuidInput)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestParseNUID_WhitespaceOnlyIsRejected(t *testing.T) {
	nuidInput := " "
	_, err := ParseNUID(nuidInput)

	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestParseNUID_EmptyStringIsRejected(t *testing.T) {
	nuidInput := ""
	_, err := ParseNUID(nuidInput)

	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestParseNUID_InvalidLengthIsRejected(t *testing.T) {
	nuidInput := strings.Repeat("1", 10)
	_, err := ParseNUID(nuidInput)

	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestParseNUID_InvalidCharactersAreRejected(t *testing.T) {
	nuidInput := "a" + strings.Repeat("1", 8)
	_, err := ParseNUID(nuidInput)

	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestParseNUID_StringWithInvalidCharactersIsRejected(t *testing.T) {
	characters := []rune{'1', 'a'}

	for numA := 1; numA <= 8; numA++ {
		permutation := strings.Repeat("a", numA)
		fullString := permutation + "11111111"[0:8-numA]

		for i := 0; i < 9; i++ {
			for _, char := range characters {
				testString := InsertCharAtIndex(fullString, char, i)
				_, err := ParseNUID(testString)

				if err == nil {
					t.Errorf("Expected an error for string: %s", testString)
				}
			}
		}
	}
}
