package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	result := pump.GlucoseTargets()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Printf("%+v\n", result)
	t := result.GlucoseTargetAt(time.Now())
	fmt.Printf("Current target: %d – %d %v\n", t.Low, t.High, t.Units)
}
