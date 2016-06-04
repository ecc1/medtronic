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
	t, err := pump.Clock()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t)
}
