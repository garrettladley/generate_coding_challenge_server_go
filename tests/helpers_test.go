package tests

import (
	"testing"

	"github.com/garrettladley/generate_coding_challenge_server_go/domain"
	"github.com/stretchr/testify/assert"
)

func TestInsertCharAtIndex(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		char     rune
		index    int
		expected string
	}{
		{"abc", 'X', 0, "Xabc"},
		{"abc", 'X', 1, "aXbc"},
		{"abc", 'X', 2, "abXc"},
		{"abc", 'X', 3, "abcX"},
		{"", 'X', 0, "X"},
		{"xyz", 'X', 3, "xyzX"},
	}

	for _, test := range tests {
		result := domain.InsertCharAtIndex(test.input, test.char, test.index)
		assert.Equal(test.expected, result)
	}
}
