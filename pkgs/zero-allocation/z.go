package main

import (
	"fmt"
	"io"
	"os"
)

const (
	bufferSize = 10
	fileName   = "textf.txt"
)

func main() {
	readChunk()
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
