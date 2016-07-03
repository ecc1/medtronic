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
	numHours = flag.Int("n", 6, "number of `hours`")
	earliest = time.Now()
)

func main() {
	flag.Parse()
	pump := medtronic.Open()
	pump.Wakeup()
	newer := pump.Family() >= 23
	numPages := pump.HistoryPageCount()
	cutoff := earliest.Add(-time.Duration(*numHours) * time.Hour)
	log.Printf("scanning for records since %v", cutoff)
	for page := 0; page < numPages && !earliest.Before(cutoff) && pump.Error() == nil; page++ {
		data := pump.HistoryPage(page)
		records, err := medtronic.DecodeHistoryRecords(data, newer)
		if err != nil {
			pump.SetError(err)
		}
		for _, r := range records {
			handleRecord(r)
		}
	}
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}

func handleRecord(r medtronic.HistoryRecord) {
	t := r.Time
	if !t.IsZero() && t.Before(earliest) {
		earliest = t
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Printf("%v %v\n", r.Type(), err)
	} else {
		fmt.Printf("%v %s\n", r.Type(), string(b))
	}
}
