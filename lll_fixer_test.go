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
		{"Hello", "Hel"},
		{"\nHello", "Hel"},
		{"abc\nxyz\nTest", "Tes"},
	}
	// loop through the tests
	for _, test := range tests {
		actual := lineLead(test.input)
		// compare the actual vs expected
		if actual != test.expected {
			t.Errorf("Test failed! Expected: %q, Actual: %q", test.expected, actual)
		}
	}
}

func TestSplitAtWord(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"123456789ABC", "123456789A\n BC"},
		{"// abc 1234567890", "// abc\n// 1234567890"},
		{"/* abc   1234567890 */", "/* abc\n * 1234567890 */"},
		{"/*\n * abc   1234567890\n*/", "/*\n * abc\n * 1234567890\n*/"},
		{`"abc 1234567890"`, `"abc" +
	" 1234567890"`},
		{`"123456789ABC"`, `"123456789" +
	"ABC"`},
	}
	// loop through the tests
	for _, test := range tests {
		actual := splitAtWord(test.input, 10)
		// compare the actual vs expected
		if actual != test.expected {
			t.Errorf("Test failed! Expected: %q, Actual: %q", test.expected, actual)
		}
	}
}
