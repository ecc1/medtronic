package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	result := pump.InsulinSensitivities()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Printf("%+v\n", result)
	s := result.InsulinSensitivityAt(time.Now())
	fmt.Printf("Current insulin sensitivity: %d\n", s.Sensitivity)
}
