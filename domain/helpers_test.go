package domain

import (
	"testing"
)

func TestInsertCharAtIndex(t *testing.T) {
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
		result := InsertCharAtIndex(test.input, test.char, test.index)
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s', but got '%s'", test.input, test.expected, result)
		}
	}
}
