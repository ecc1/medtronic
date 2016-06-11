package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/medtronic"
)

func usage() {
	log.Fatalf("Usage: %s duration units/hr", os.Args[0])
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}
	duration, err := time.ParseDuration(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	rate, err := strconv.ParseFloat(os.Args[2], 32)
	if err != nil {
		log.Fatal(err)
	}
	pump := medtronic.Open()
	log.Printf("setting temporary basal of %.3f units/hour for %v", rate, duration)
	pump.SetTempBasal(duration, int(rate*1000))
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}

func parseTime(date string) time.Time {
	const layout = "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, date, time.Local)
	if err != nil {
		log.Fatalf("Cannot parse %s: %v", date, err)
	}
	return t
}
