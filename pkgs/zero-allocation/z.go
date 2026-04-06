package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	bufferSize = 10
	fileName   = "textf.txt"
)

func main() {
	var r io.Reader = strings.NewReader("some io.Reader stream to be read\n")

	myTee := &myTeeReader{r: r, w: os.Stdout}
	if _, err := io.ReadAll(myTee); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from myTeeReader:", err)
	}
}

type myTeeReader struct {
	r io.Reader
	w io.Writer
}

func (t *myTeeReader) Read(p []byte) (n int, err error) {
	till, err := t.r.Read(p)
	if till > 0 {
		if _, err := t.w.Write(p[:till]); err != nil {
			return till, err
		}
	}
	return till, err
}

func readChunk() {
	buffer := make([]byte, bufferSize)
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	defer file.Close()

	if err != nil {
		os.Exit(1)
	}

	for {
		b, err := file.Read(buffer)
		if b > 0 {
			fmt.Printf("%s", buffer[:b])
			continue
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			os.Exit(1)
		}
	}
}
