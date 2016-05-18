package main

import (
	"fmt"
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
	id, err := dev.PumpID(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}
