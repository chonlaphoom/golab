package main

import (
	greeter "github.com/chonlaphoom/golab/greeting"
	subgreeter "github.com/chonlaphoom/golab/greeting/subgreeting"
	"github.com/chonlaphoom/golab/rune"
	"log"
)

func main() {
	log.Println("Running Rune Example:")
	rune.Print()

	log.Println("Running Greeting Example:")
	greeter.Greeting()
	subgreeter.SubGreeting()
}
