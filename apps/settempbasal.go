package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/medtronic"
)

func usage() {
	log.Fatalf("Usage: %s duration (units/hr | rate%%)", os.Args[0])
}

func main() {
	if len(os.Args) != 3 || len(os.Args[1]) == 0 || len(os.Args[2]) == 0 {
		usage()
	}
	duration, err := time.ParseDuration(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	pump := medtronic.Open()
	rateArg := os.Args[2]
	n := len(rateArg) - 1
	if rateArg[n] == '%' {
		percent, err := strconv.Atoi(rateArg[:n])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("setting temporary basal of %d%% for %v", percent, duration)
		pump.SetPercentTempBasal(duration, percent)
	} else {
		f, err := strconv.ParseFloat(rateArg, 32)
		if err != nil {
			log.Fatal(err)
		}
		rate := int(1000.0*f + 0.5)
		log.Printf("setting temporary basal of %d.%03d units/hour for %v", rate/1000, rate%1000, duration)
		pump.SetAbsoluteTempBasal(duration, rate)
	}
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}
