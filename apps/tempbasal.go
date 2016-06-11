package main

import (
	"fmt"
	"log"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	result := pump.TempBasal()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Printf("%+v\n", result)
}
