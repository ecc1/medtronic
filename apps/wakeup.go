package main

import (
	"log"

	"github.com/ecc1/cc1100"
)

func main() {
	dev, err := cc1100.Open()
	if err != nil {
		log.Fatal(err)
	}
	err = dev.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = dev.Wakeup()
	if err != nil {
		log.Fatal(err)
	}
}
