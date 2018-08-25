package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
)

var (
	all       = flag.Bool("a", false, "get entire pump history")
	numHours  = flag.Int("n", 6, "number of `hours` of history to get")
	nsFlag    = flag.Bool("ns", false, "format as Nightscout treatments")
	sinceFlag = flag.String("s", "", "get history since the specified `time` in RFC3339 format")
	fromFlag = flag.String("f", "", "get history from a specified `record id` where the record id the base64 encoded binary data of the given record")
)

func main() {
	flag.Parse()
	var cutoff time.Time
	var err error
	var from []byte = nil;
	if *all {
		log.Printf("retrieving entire pump history")
	} else if *fromFlag != "" {
		from, err = base64.StdEncoding.DecodeString(*fromFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else if *sinceFlag != "" {
		cutoff, err = time.Parse(medtronic.JSONTimeLayout, *sinceFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cutoff = time.Now().Add(-time.Duration(*numHours) * time.Hour)
	}
	if from != nil {
		log.Printf("retrieving pump history from %s", base64.StdEncoding.EncodeToString(from))
	} else if !*all {
		log.Printf("retrieving pump history since %s", cutoff.Format(medtronic.UserTimeLayout))
	}
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	var results medtronic.History
	if from != nil {
		var found bool
		found, results = pump.HistoryFrom(from)
		if !found {
			if pump.Error() != nil {
				log.Fatal(pump.Error())
			} else {
				os.Exit(2)
			}
		}
	} else {
		results = pump.History(cutoff)
	}
	if *nsFlag {
		medtronic.ReverseHistory(results)
		fmt.Println(nightscout.JSON(medtronic.Treatments(results)))
	} else {
		fmt.Println(nightscout.JSON(results))
	}
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}
