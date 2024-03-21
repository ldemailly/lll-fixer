package main

import (
	"testing"
)

func TestLineLead(t *testing.T) {
	// tests input vs actual struct:
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"\nHello", "Hello"},
		{"abc\nxyz\nTest", "Test"},
	}
	// loop through the tests
	for _, test := range tests {
		actual := lineLead(test.input)
		// compare the actual vs expected
		if actual != test.expected {
			t.Errorf("Test failed! Expected: %v, Actual: %v", test.expected, actual)
		}
	}
}
