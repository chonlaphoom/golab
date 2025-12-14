package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	defer func() {
		writer.Flush()
		file.Close()
	}()

	joinStrings := func(items []string) string {
		return strings.Join(items, ";")
	}

	// Write CSV header
	err = writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	if err != nil {
		return err
	}
	for _, pageData := range pages {
		err := writer.Write([]string{
			pageData.URL,
			pageData.H1,
			pageData.FirstParagraph,
			joinStrings(pageData.OutgoingLinks),
			joinStrings(pageData.ImageURLs),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
