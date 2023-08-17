package domain

import (
	"strings"
	"testing"
)

func TestParseApplicantName_ValidName(t *testing.T) {
	nameInput := "Muneer Lalji"
	_, err := ParseApplicantName(nameInput)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestParseApplicantName_256GraphemeLongNameIsValid(t *testing.T) {
	nameInput := strings.Repeat("a", 256)
	_, err := ParseApplicantName(nameInput)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestParseApplicantName_NameLongerThan256GraphemesIsRejected(t *testing.T) {
	nameInput := strings.Repeat("a", 257)
	_, err := ParseApplicantName(nameInput)

	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestParseApplicantName_WhitespaceOnlyNamesAreRejected(t *testing.T) {
	nameInput := " "
	_, err := ParseApplicantName(nameInput)

	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestParseApplicantName_EmptyStringIsRejected(t *testing.T) {
	nameInput := ""
	_, err := ParseApplicantName(nameInput)

	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestParseApplicantName_NamesContainingInvalidCharactersAreRejected(t *testing.T) {
	invalidCharacters := []string{"/", "(", ")", "\"", "<", ">", "\\", "{", "}"}

	for _, char := range invalidCharacters {
		nameInput := char
		_, err := ParseApplicantName(nameInput)

		if err == nil {
			t.Errorf("Expected an error for character %s, but got none", char)
		}
	}
}
