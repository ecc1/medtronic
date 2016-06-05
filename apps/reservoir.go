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
	result, err := pump.Reservoir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%.2f\n", float64(result)/1000)
}
