package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	result := pump.BasalRates()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Printf("%+v\n", result)
	rate := result.BasalRateAt(time.Now()).Rate
	fmt.Printf("Current rate: %d.%03d units/hour\n", rate/1000, rate%1000)
}
