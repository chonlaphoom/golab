package main

import (
	"net/url"
)

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func normalizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	// remove trailing slash if present
	normalizedPath := parsedURL.Path
	if len(normalizedPath) > 0 && normalizedPath[len(normalizedPath)-1] == '/' {
		normalizedPath = normalizedPath[:len(normalizedPath)-1]
	}

	normalizedDomainAndPath := parsedURL.Host + normalizedPath
	return normalizedDomainAndPath, nil
}
