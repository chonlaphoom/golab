package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func newConfig(baseURL *url.URL, maxConcurrency, maxPages int) *config {
	return &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, found := cfg.pages[normalizedURL]
	return !found
}

func (cfg *config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = data
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		// fmt.Printf("Finished processing: %s\n", rawCurrentURL)
		cfg.wg.Done()
	}()

	if len(cfg.pages) >= cfg.maxPages {
		// fmt.Printf("Reached max page limit of %d, stopping crawl.\n", cfg.maxPages)
		return
	}

	parsedRawBaseURL := cfg.baseURL
	parsedRawCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if parsedRawBaseURL.Hostname() != parsedRawCurrentURL.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if isFirst := cfg.addPageVisit(normalizedCurrentURL); !isFirst {
		return
	}

	rawHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}
	fmt.Printf("[%s] Crawled: %s\n", time.Now().Format(time.RFC3339), rawCurrentURL)

	pageData := extractPageData(rawHTML, rawCurrentURL)
	cfg.setPageData(normalizedCurrentURL, pageData)

	for _, link := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(link)
	}

	time.Sleep(500 * time.Millisecond) // polite delay between requests
}
