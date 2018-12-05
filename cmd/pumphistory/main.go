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

const (
	success        = 0
	commError      = 1
	recordNotFound = 2
	partialResult  = 3
)

var (
	all       = flag.Bool("a", false, "get entire pump history")
	numHours  = flag.Int("n", 6, "number of `hours` of history to get")
	nsFlag    = flag.Bool("ns", false, "format as Nightscout treatments")
	fromFlag  = flag.String("f", "", "get history since the specified record `ID` (the base64-encoding the record data)")
	sinceFlag = flag.String("s", "", "get history since the specified `time` in RFC3339 format")

	cutoff   time.Time
	recordID []byte
)

func main() {
	parseFlags()
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	var results medtronic.History
	found := true
	if *fromFlag != "" {
		results, found = pump.HistoryFrom(recordID)
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
		log.Print(pump.Error())
		if len(results) != 0 {
			os.Exit(partialResult)
		}
		os.Exit(commError)
	}
	if !found {
		log.Printf("record %s not found", *fromFlag)
		os.Exit(recordNotFound)
	}
	os.Exit(success)
}

func parseFlags() {
	flag.Parse()
	var err error
	if *all {
		log.Printf("retrieving entire pump history")
	} else if *fromFlag != "" {
		recordID, err = base64.StdEncoding.DecodeString(*fromFlag)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("retrieving pump history since record %s", *fromFlag)
	} else if *sinceFlag != "" {
		cutoff, err = time.Parse(medtronic.JSONTimeLayout, *sinceFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cutoff = time.Now().Add(-time.Duration(*numHours) * time.Hour)
	}
	if !*all && *fromFlag == "" {
		log.Printf("retrieving pump history since %s", cutoff.Format(medtronic.UserTimeLayout))
	}
}
