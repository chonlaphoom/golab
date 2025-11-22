package main

import (
	"golab/internal/rune"
	greeter "golab/pkgs/greeting"
	subgreeter "golab/pkgs/greeting/subgreeting"
	"log"
)

func main() {
	log.Println("Running Rune Example:")
	rune.Print()

	log.Println("Running Greeting Example:")
	greeter.Greeting()
	subgreeter.SubGreeting()
}
