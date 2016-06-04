package main

import (
	"log"
	"os"
	"time"

	"github.com/ecc1/medtronic"
)

func usage() {
	log.Fatalf("Usage: %s YYYY-MM-DD HH:MM:SS (or \"now\")\n", os.Args[0])
}

func main() {
	var t time.Time
	switch len(os.Args) {
	case 3:
		t = parseTime(os.Args[1] + " " + os.Args[2])
	case 2:
		if os.Args[1] == "now" {
			t = time.Now()
		} else {
			usage()
		}
	default:
		usage()
	}
	pump, err := medtronic.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("setting pump clock to %s\n", t.Format(time.UnixDate))
	err = pump.SetClock(t)
	if err != nil {
		log.Fatal(err)
	}
}

func parseTime(date string) time.Time {
	const layout = "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, date, time.Local)
	if err != nil {
		log.Fatalf("Cannot parse %s: %v\n", date, err)
	}
	return t
}
