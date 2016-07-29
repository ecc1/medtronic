package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

var (
	all      = flag.Bool("a", false, "get entire pump history")
	numHours = flag.Int("n", 6, "number of `hours` of history to get")
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
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
