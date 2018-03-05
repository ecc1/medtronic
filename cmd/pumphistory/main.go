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
	all       = flag.Bool("a", false, "get entire pump history")
	numHours  = flag.Int("n", 6, "number of `hours` of history to get")
	nsFlag    = flag.Bool("t", false, "format as Nightscout treatments")
	sinceFlag = flag.String("s", "", "get history since the specified `time` in RFC3339 format")
)

func main() {
	flag.Parse()
	var cutoff time.Time
	var err error
	if *all {
		log.Printf("retrieving entire pump history")
	} else if *sinceFlag != "" {
		cutoff, err = time.Parse(medtronic.JSONTimeLayout, *sinceFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cutoff = time.Now().Add(-time.Duration(*numHours) * time.Hour)
	}
	if !*all {
		log.Printf("retrieving pump history since %s", cutoff.Format(medtronic.UserTimeLayout))
	}
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	results := pump.History(cutoff)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	if *nsFlag {
		medtronic.ReverseHistory(results)
		fmt.Println(nightscout.JSON(medtronic.Treatments(results)))
	} else {
		fmt.Println(nightscout.JSON(results))
	}
}
