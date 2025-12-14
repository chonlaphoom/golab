package main

import (
	"net/url"
	"reflect"
	"testing"
)

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

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://blog.boot.dev"

	_, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "absolute URL",
			inputURL:  inputURL,
			inputBody: `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://blog.boot.dev"},
		}, {
			name:      "relative URL",
			inputURL:  inputURL,
			inputBody: `<html><body><a href="/"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://blog.boot.dev"},
		}, {
			name:     "mixed URLs",
			inputURL: inputURL,
			inputBody: `<html><body>
			<a href="https://blog.boot.dev"><span>Boot.dev</span></a>
			<a href="/about"><span>About</span></a>
		</body></html>`,
			expected: []string{"https://blog.boot.dev", "https://blog.boot.dev/about"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, inputURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}

}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	_, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "relative image URL",
			input:    `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected: []string{"https://blog.boot.dev/logo.png"},
		},
		{
			name:     "absolute image URL",
			input:    `<html><body><img src="https://blog.boot.dev/logo.png" alt="Logo"></body></html>`,
			expected: []string{"https://blog.boot.dev/logo.png"},
		},
		{
			name:     "mixed image URLs",
			input:    `<html><body><img src="/logo.png" alt="Logo"><img src="https://blog.boot.dev/banner.png" alt="Banner"></body></html>`,
			expected: []string{"https://blog.boot.dev/logo.png", "https://blog.boot.dev/banner.png"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getImagesFromHTML(tc.input, inputURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
