package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/ecc1/medtronic"
)

var (
	model = flag.Int("m", 523, "pump model")

	timeBlank = strings.Repeat(" ", len(medtronic.UserTimeLayout))
)

func main() {
	flag.Parse()
	newer := *model%100 > 22
	for _, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		d := json.NewDecoder(f)
		var records []medtronic.HistoryRecord
		err = d.Decode(&records)
		if err != nil {
			log.Fatal(err)
		}
		f.Close() // nolint
		for _, r := range records {
			validate(r, newer)
		}
	}
}

func validate(r medtronic.HistoryRecord, newerPump bool) {
	rr, err := medtronic.DecodeHistoryRecord(r.Data, newerPump)
	t := time.Time(r.Time)
	tStr := timeBlank
	if !t.IsZero() {
		tStr = t.Format(medtronic.UserTimeLayout)
	}
	if err != nil {
		log.Printf("%s %v: %v\n", tStr, r.Type(), err)
	}
	if reflect.DeepEqual(r, rr) {
		fmt.Printf("%s %v\n", tStr, r.Type())
	} else {
		fmt.Printf("Records do not match:\n")
		fmt.Printf("  %+v\n", r)
		fmt.Printf("  %+v\n", rr)
	}
}
