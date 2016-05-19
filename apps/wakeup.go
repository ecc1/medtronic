package main

import (
	"log"

	"github.com/ecc1/medtronic"
)

func main() {
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	err = pump.Wakeup()
	if err != nil {
		log.Fatal(err)
	}
}
