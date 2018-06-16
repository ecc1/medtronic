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
	all       = flag.Bool("a", false, "get entire CGM history")
	numHours  = flag.Int("n", 6, "number of `hours` of history to get")
	nsFlag    = flag.Bool("e", false, "format as Nightscout entries")
	sinceFlag = flag.String("s", "", "get history since the specified `time` in RFC3339 format")
	noTimes   = flag.Bool("notimes", false, "do not add times to glucose records")
)

func main() {
	flag.Parse()
	var cutoff time.Time
	var err error
	if *all {
		log.Printf("retrieving entire CGM history")
	} else if *sinceFlag != "" {
		cutoff, err = time.Parse(medtronic.JSONTimeLayout, *sinceFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cutoff = time.Now().Add(-time.Duration(*numHours) * time.Hour)
	}
	if !*all {
		log.Printf("retrieving CGM history since %s", cutoff.Format(medtronic.UserTimeLayout))
	}
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	results := pump.CGMHistory(cutoff)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	if !*noTimes {
		medtronic.AddCGMTimes(results)
	}
	if *nsFlag {
		medtronic.ReverseCGMHistory(results)
		fmt.Println(nightscout.JSON(medtronic.Entries(results)))
	} else {
		fmt.Println(nightscout.JSON(results))
	}
}
