package main

import (
	"fmt"
	"net/url"
	"time"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) map[string]int {
	parsedRawBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return pages
	}

	parsedRawCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return pages
	}

	// only crawl pages on the same domain
	if parsedRawBaseURL.Host != parsedRawCurrentURL.Host {
		return pages
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return pages
	}

	if _, ok := pages[normalizedCurrentURL]; ok {
		pages[normalizedCurrentURL]++
		return pages
	}

	// mark page as visited
	pages[normalizedCurrentURL] = 1

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return pages
	}
	fmt.Printf("crawled page: %s\n", rawCurrentURL)
	urls, err := getURLsFromHTML(html, rawCurrentURL)
	if err != nil {
		return pages
	}

	// delay between requests to avoid overwhelming the server
	time.Sleep(3 * time.Second)

	for _, link := range urls {
		fmt.Printf("found link: %s\n", link)
		crawlPage(rawBaseURL, link, pages)
	}

	return pages
}
