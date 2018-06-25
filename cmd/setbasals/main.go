package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ecc1/medtronic"
)

var (
	hourly   = flag.Bool("h", false, "assume hourly times starting at 00:00")
	simulate = flag.Bool("s", false, "print schedule only; do not send to pump")

	pump *medtronic.Pump
)

func main() {
	flag.Parse()
	var sched medtronic.BasalRateSchedule
	if *hourly {
		sched = hourlySchedule(flag.Args())
	} else {
		sched = arbitrarySchedule(flag.Args())
	}
	if len(sched) == 0 {
		log.Fatal("cannot set an empty schedule")
	}
	if sched[0].Start != 0 {
		log.Fatal("schedule must begin at 00:00")
	}
	if *simulate {
		showJSON(sched)
		return
	}
	pump = medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	pump.SetBasalRates(sched)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
}

func hourlySchedule(rates []string) medtronic.BasalRateSchedule {
	var sched medtronic.BasalRateSchedule
	for i, arg := range rates {
		d := time.Duration(i) * time.Hour
		br := medtronic.BasalRate{
			Start: medtronic.Duration(d).TimeOfDay(),
			Rate:  getRate(arg),
		}
		sched = append(sched, br)

	}
	return sched
}

func arbitrarySchedule(args []string) medtronic.BasalRateSchedule {
	var sched medtronic.BasalRateSchedule
	if len(args)%2 != 0 {
		log.Fatal("schedule must be specified as a sequence of start times and rates")
	}
	for i := 0; i < len(args)-1; i += 2 {
		br := medtronic.BasalRate{
			Start: parseTD(args[i]),
			Rate:  getRate(args[i+1]),
		}
		sched = append(sched, br)
	}
	return sched
}

func getRate(arg string) medtronic.Insulin {
	f, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		log.Fatal(err)
	}
	return medtronic.Insulin(1000.0*f + 0.5)
}

func parseTD(s string) medtronic.TimeOfDay {
	td, err := medtronic.ParseTimeOfDay(s)
	if err != nil {
		log.Fatal(err)
	}
	return td
}

func showJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
		fmt.Println(v)
		return
	}
	fmt.Println(string(b))
}
