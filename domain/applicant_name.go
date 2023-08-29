package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ApplicantName string

func ParseApplicantName(str string) (*ApplicantName, error) {
	isEmptyOrWhitespace := strings.TrimSpace(str) == ""
	isTooLong := utf8.RuneCountInString(str) > 256

	forbiddenCharacters := []rune{'/', '(', ')', '"', '<', '>', '\\', '{', '}'}
	containsForbiddenCharacters := false

	for _, char := range str {
		if containsRune(forbiddenCharacters, char) {
			containsForbiddenCharacters = true
			break
		}
	}

	if isEmptyOrWhitespace || isTooLong || containsForbiddenCharacters {
		return nil, fmt.Errorf("invalid name! Given: %s", str)
	}

	applicantName := ApplicantName(str)
	return &applicantName, nil
}

func (name *ApplicantName) String() string {
	return string(*name)
}

func containsRune(slice []rune, target rune) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}
