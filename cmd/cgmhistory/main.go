package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"os"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/nightscout"
)

var (
	all       = flag.Bool("a", false, "get entire CGM history")
	numHours  = flag.Int("n", 6, "number of `hours` of history to get")
	nsFlag    = flag.Bool("e", false, "force output in stdout and format as Nightscout entries")
	sinceFlag = flag.String("s", "", "get history since the specified `time` in RFC3339 format")
	jsonFile  = flag.String("l", "", "append results to `file` in legacy format")
	nsFile    = flag.String("ns", "", "append results to `file` in Nightscout format")
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
		log.Printf("retrieving pump history since %s", cutoff.Format(medtronic.UserTimeLayout))
	}
	
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	results := pump.CGMHistory(cutoff)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	
	if *nsFlag {
		medtronic.ReverseCGMHistory(results)
		fmt.Println(nightscout.JSON(medtronic.Entries(results)))
	} else if *jsonFile == "" && *nsFile == "" {
		fmt.Println(nightscout.JSON(results))
	}
	
	if *jsonFile != "" {
		formattedData := medtronic.FormatToOAPS(results)
		jsonData := nightscout.JSON(formattedData)
		f, err := os.Create(*jsonFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.WriteString(jsonData)
		log.Printf("wrote %d entries to %s", len(formattedData), *jsonFile)
	}
	
	if *nsFile != "" {
		medtronic.ReverseCGMHistory(results)
		enties := medtronic.Entries(results)
		jsonData := nightscout.JSON(enties)
		f, err := os.Create(*nsFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.WriteString(jsonData)
		log.Printf("wrote %d entries to %s", len(enties), *nsFile)
	}
}
