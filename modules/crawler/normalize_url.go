package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func getURLsFromHTML(htmlBody string, baseURL string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		val, exist := item.Attr("href")
		if !exist {
			return
		}

		// is href value an absolute URL?
		hrefURL, errorParsing := url.Parse(val)
		if errorParsing != nil {
			return
		}

		if !hrefURL.IsAbs() {
			if hrefURL.Path == "/" {
				urls = append(urls, baseURL)
				return
			}

			base, err := hrefURL.Parse(baseURL)
			if err != nil {
				return
			}
			hrefURL = base.ResolveReference(hrefURL)
			urls = append(urls, hrefURL.String())
			return
		}

		urls = append(urls, val)
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var imageURLs []string
	doc.Find("img[src]").Each(func(_ int, item *goquery.Selection) {
		val, exist := item.Attr("src")
		if !exist {
			return
		}

		// is src value an absolute URL?
		srcURL, errorParsing := url.Parse(val)
		if errorParsing != nil {
			return
		}

		if !srcURL.IsAbs() {
			if srcURL.Path == "/" {
				imageURLs = append(imageURLs, baseURL)
				return
			}

			base, err := srcURL.Parse(baseURL)
			if err != nil {
				return
			}
			srcURL = base.ResolveReference(srcURL)
			imageURLs = append(imageURLs, srcURL.String())
			return
		}

		imageURLs = append(imageURLs, val)
	})

	return imageURLs, nil
}

func extractPageData(htmlBody, pageURL string) PageData {
	h1 := getH1FromHTML(htmlBody)
	firstParagraph := getFirstParagraphFromHTML(htmlBody)

	outgoingLinks, err := getURLsFromHTML(htmlBody, pageURL)
	if err != nil {
		return PageData{}
	}

	imageURLs, err := getImagesFromHTML(htmlBody, pageURL)
	if err != nil {
		return PageData{}
	}

	return PageData{
		URL:            pageURL,
		H1:             h1,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
