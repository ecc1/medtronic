package main

import (
	"log"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	pump.Wakeup()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}
