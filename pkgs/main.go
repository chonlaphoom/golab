package main

import (
	greeter "github.com/chonlaphoom/golab/pkgs/greeting"
	subgreeter "github.com/chonlaphoom/golab/pkgs/greeting/subgreeting"
	"github.com/chonlaphoom/golab/pkgs/rune"
	"github.com/chonlaphoom/golab/pkgs/throttler"
	"log"
)

func main() {
	log.Println("Running Rune Example:")
	rune.Print()

	log.Println("Running Greeting Example:")
	greeter.Greeting()
	subgreeter.SubGreeting()

	log.Println("Running Throttler Example:")
	throttler.Execute()
}
