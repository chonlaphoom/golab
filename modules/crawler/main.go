package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	actualArgs := os.Args[1:]
	if len(actualArgs) < 3 {
		fmt.Println("no website provided")
		os.Exit(1)
		return
	}

	if len(actualArgs) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
		return
	}

	baseURL := actualArgs[0]
	maxConcurrencyStr := actualArgs[1]
	maxPagesStr := actualArgs[2]

	maxConcurrent, err := strconv.Atoi(maxConcurrencyStr)
	if err != nil || maxConcurrent <= 0 {
		fmt.Println("invalid max concurrency value")
		os.Exit(1)
		return
	}
	if maxConcurrent > 5 {
		fmt.Println("max concurrency value is 5")
		maxConcurrent = 5
	}

	maxPages, err := strconv.Atoi(maxPagesStr)
	if err != nil || maxPages <= 0 {
		fmt.Println("invalid max pages value")
		os.Exit(1)
		return
	}

	// fmt.Printf("Max Concurrency: %d\n", maxConcurrent)
	// fmt.Printf("Max Pages: %d\n", maxPages)
	fmt.Printf("starting crawl\n%s\n\n", baseURL)

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("error parsing base URL: %v\n", err)
		os.Exit(1)
		return
	}

	cfg := newConfig(parsedBaseURL, maxConcurrent, maxPages)

	cfg.wg.Add(1)
	cfg.crawlPage(baseURL)

	cfg.wg.Wait()
	fmt.Printf("crawl finished\n")
	for _, pageData := range cfg.pages {
		fmt.Printf("Found page: %s\n", pageData.URL)
	}

	writeCSVReport(cfg.pages, "report.csv")
	fmt.Printf("\nreport generated: report.csv\n")
	os.Exit(0)
}

func getHTML(rawURL string) (string, error) {
	request, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("User-Agent", "MyCrawler/1.0") // Set a custom User-Agent
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode > 400 {
		return "", fmt.Errorf("received status code %d", response.StatusCode)
	}

	if !strings.Contains(
		response.Header.Get("Content-Type"),
		"text/html",
	) {
		return "", fmt.Errorf("invalid content type: %s", response.Header.Get("Content-Type"))
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
