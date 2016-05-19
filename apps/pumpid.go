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
	id, err := pump.ID(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}
