package main

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
