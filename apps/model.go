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
	model, err := pump.Model()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(model)
}
