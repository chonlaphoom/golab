package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "with https",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "with http",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "start with http and ends with slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "start with https and ends with slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
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

func TestGetH1FromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetEmptyFirstParagraphFromHTML(t *testing.T) {
	inputBody := `<html><body>
		<div>No paragraphs here!</div>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
