package main

import (
	"fmt"
	"log"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	result := pump.Reservoir()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Printf("%.2f\n", float64(result)/1000)
}
