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
	buffer := make([]byte, bufferSize)
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	defer file.Close()

	if err != nil {
		os.Exit(1)
	}

	for {
		b, err := file.Read(buffer)
		if b == 0 || err == io.EOF {
			break
		}
		if err != nil {
			os.Stderr.WriteString("Error reading file: " + err.Error() + "\n")
			break
		}
		fmt.Printf("%s", buffer[:b])
	}
}
