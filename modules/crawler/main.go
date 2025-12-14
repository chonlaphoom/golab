package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	actualArgs := os.Args[1:]
	if len(actualArgs) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
		return
	}

	if len(actualArgs) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
		return
	}

	fmt.Printf("starting crawl\n%s\n", actualArgs[0])
	html, err := getHTML(actualArgs[0])
	if err != nil {
		fmt.Printf("error fetching HTML: %s\n", err.Error())
		os.Exit(1)
		return
	}

	fmt.Printf("fetched HTML\n%s\n", html)

	fmt.Println("crawl finished")
	os.Exit(0)
}

func getHTML(rawURL string) (string, error) {
	request, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("User-Agent", "BootCrawler/1.0")
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
