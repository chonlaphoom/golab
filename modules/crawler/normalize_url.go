package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func getH1FromHTML(htmlBody string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return ""
	}
	return doc.Find("h1").First().Text()
}

func getFirstParagraphFromHTML(htmlBody string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return ""
	}

	mainSelection := doc.Find("main")
	foundMain := mainSelection.Length() > 0
	if foundMain {
		return mainSelection.Find("p").First().Text()
	} else {
		return doc.Find("p").First().Text()
	}
}
