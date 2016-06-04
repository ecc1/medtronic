package main

import (
	"fmt"
	"log"

	"github.com/ecc1/medtronic"
)

func main() {
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	status, err := pump.PumpStatus(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", status)
}
