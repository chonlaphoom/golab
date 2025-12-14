package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "with https",
			inputURL: "https://blog.test.dev/path",
			expected: "blog.test.dev/path",
		},
		{
			name:     "with http",
			inputURL: "http://blog.test.dev/path",
			expected: "blog.test.dev/path",
		},
		{
			name:     "start with http and ends with slash",
			inputURL: "http://blog.test.dev/path/",
			expected: "blog.test.dev/path",
		},
		{
			name:     "start with https and ends with slash",
			inputURL: "https://blog.test.dev/path/",
			expected: "blog.test.dev/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
