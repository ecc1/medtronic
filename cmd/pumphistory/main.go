package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
)

var (
	all      = flag.Bool("a", false, "get entire pump history")
	numHours = flag.Int("n", 6, "number of `hours` of history to get")
	nsFlag   = flag.Bool("t", false, "format as Nightscout treatments")
)

func main() {
	flag.Parse()
	cutoff := time.Time{}
	if *all {
		log.Printf("retrieving entire pump history")
	} else {
		cutoff = time.Now().Add(-time.Duration(*numHours) * time.Hour)
		log.Printf("retrieving pump history since %s", cutoff.Format(medtronic.UserTimeLayout))
	}
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	results := pump.HistoryRecords(cutoff)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	if *nsFlag {
		medtronic.ReverseHistoryRecords(results)
		fmt.Println(nightscout.JSON(medtronic.Treatments(results)))
	} else {
		fmt.Println(nightscout.JSON(results))
	}
}
